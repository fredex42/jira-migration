package trello

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fredex42/mm-jira-migration/common"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadTrelloAttachment(cardId string, fileName string, jiraAttachment *common.Attachment, apiKey *common.ScriptKey, httpClient *http.Client) error {
	uri := fmt.Sprintf("https://api.trello.com/1/cards/%s/attachments?key=%s&token=%s", cardId, apiKey.User, apiKey.Key)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", jiraAttachment.Filename)
	io.Copy(part, file)

	mt, _ := writer.CreateFormField("mimeType")
	mt.Write([]byte(jiraAttachment.MimeType))

	writer.Close()

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	response, err := httpClient.Do(req)

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		log.Printf("ERROR Server responded %s", string(responseContent))
		return errors.New(fmt.Sprintf("could not create attachment, server responded with a %d", response.StatusCode))
	}

	log.Printf("INFO Uploaded attachment from %s to Trello for %s", fileName, jiraAttachment.Filename)
	os.Remove(fileName)
	return nil
}
