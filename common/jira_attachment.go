package common

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

/*
DownloadJiraAttachment downloads the content of the given attachment to a temp file.
Returns the name of the downloaded file on success, or an error.
*/
func DownloadJiraAttachment(attachmentId string, hostname *string, apiKey *ScriptKey, httpClient *http.Client) (string, error) {
	uri := fmt.Sprintf("https://%s/rest/api/3/attachment/content/%s?redirect=false", *hostname, attachmentId)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(apiKey.User, apiKey.Key)

	response, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		io.Copy(ioutil.Discard, response.Body)
		return "", errors.New(fmt.Sprintf("server returned %d", response.StatusCode))
	}

	file, err := ioutil.TempFile("/tmp", "trelloatt")
	if err != nil {
		return "", err
	}
	defer file.Close() //we re-open for reading

	bytesCopied, err := io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	log.Printf("INFO Downloaded %d bytes of attachment to '%s'", bytesCopied, file.Name())
	return file.Name(), nil
}
