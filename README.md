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
	var config backblaze.DownloadUrlTokenParams

  // get data from POST req
	data := json.NewDecoder(req.Body)
	err := data.Decode(&token)
  if err != nil {
    panic(err)
  }


  // log into backblaze with accountId, and ApplicationKey
	user, err := backblaze.GetUser("your-backblaze-key-here")
  if err != nil {
    panic(err)
  }

  // get file auth url
	downloadUrl, err := user.GetFileDownloadUrl(token)
  if err != nil {
    panic(err)
  }
}

```
