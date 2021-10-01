package backblaze

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type user struct {
	AbsoluteMinimumPartSize int    `json:"absoluteMinimumPartSize"`
	AccountID               string `json:"accountId"`
	Allowed                 struct {
		BucketID     string      `json:"bucketId"`
		BucketName   string      `json:"bucketName"`
		Capabilities []string    `json:"capabilities"`
		NamePrefix   interface{} `json:"namePrefix"`
	} `json:"allowed"`
	APIURL              string `json:"apiUrl"`
	AuthorizationToken  string `json:"authorizationToken"`
	DownloadURL         string `json:"downloadUrl"`
	RecommendedPartSize int    `json:"recommendedPartSize"`
	S3APIURL            string `json:"s3ApiUrl"`
}

type DownloadUrlTokenParams struct {
	BucketID   string `json:"bucketId"`
	FileName   string `json:"fileName"`
	BucketName string `json:"bucketName"`
	Duration   string `json:"duration"`
}

// log into backblaze with accountId, and ApplicationKey
// string should be passed as <applicationKeyId:applicationKey>
// refer to documentation https://www.backblaze.com/b2/docs/b2_authorize_account.html
func GetUser(authKey string) (*user, error) {
	client := http.Client{}
	authorizeAccountURL := "https://api.backblazeb2.com/b2api/v2/b2_authorize_account"

	req, err := http.NewRequest("GET", authorizeAccountURL, nil)
	var userInfo user

	if err != nil {
		return nil, errors.New("error creating request")
	}
	b64Key := fmt.Sprint("Basic ", base64.StdEncoding.EncodeToString([]byte(os.Getenv(authKey))))
	req.Header = http.Header{
		"Authorization": []string{b64Key},
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("http error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("error converting response to bytes")
	}

	err = json.Unmarshal([]byte(body), &userInfo)
	if err != nil {
		return nil, errors.New("Count not unmarshal json")
	}
	return &userInfo, nil
}

// Returns download url given a users Authorization Token, and paramaters
func (u *user) GetFileDownloadUrl(token DownloadUrlTokenParams) (string, error) {

	type respModel struct {
		AuthorizationToken string `json:"authorizationToken"`
		BucketID           string `json:"bucketId"`
	}

	var temporaryDownloadToken respModel

	url := fmt.Sprintf(
		"https://api002.backblazeb2.com/b2api/v2/b2_get_download_authorization?bucketId=" + token.BucketID +
			"&fileNamePrefix=" + token.FileName + "&validDurationInSeconds=" + token.Duration,
	)

	// Get temporaryDownloadToken for single file ( or possibly group )
	c := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Authorization": []string{u.AuthorizationToken},
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &temporaryDownloadToken)
	if err != nil {
		return "", err
	}

	// piece link from authToken and fileName
	completeDownloadUrl := "https://f002.backblazeb2.com/file/" + token.BucketName + "/" + token.FileName + "?Authorization=" + temporaryDownloadToken.AuthorizationToken

	return completeDownloadUrl, nil
}
