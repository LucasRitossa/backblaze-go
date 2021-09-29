# backblaze-go
Go client library for backblaze B2 API

## Usage

Login, create download link as POST

```go
import (
	"github.com/backblaze-go"
  "net/http"
  "encoding/json"
)

func GetLink(w http.ResponseWriter, req *http.Request) {
  // get data from POST req
	data := json.NewDecoder(req.Body)
	err := data.Decode(&token)

	var token backblaze.TokenParams

  // log into backblaze account with key
	user, _ := backblaze.GetDownloadUser("your-backblaze-key-here")

  // get file auth url
	downloadUrl, _ := backblaze.GetFileDownloadUrl(user.AuthorizationToken, token)
}

```
