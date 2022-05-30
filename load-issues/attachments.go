package main

import (
	"github.com/fredex42/mm-jira-migration/common"
	"github.com/fredex42/mm-jira-migration/trello"
	"log"
	"net/http"
)

func HandleAttachments(attachmentList *[]common.Attachment, cardId string, hostname *string, jiraKey *common.ScriptKey, trelloKey *common.ScriptKey, httpClient *http.Client) error {
	log.Printf("INFO Got %d attachments", len(*attachmentList))

	for _, a := range *attachmentList {
		downloadedFileName, err := common.DownloadJiraAttachment(a.Id, hostname, jiraKey, httpClient)
		if err != nil {
			log.Printf("ERROR Could not download %s: %s", a.Filename, err)
			return err
		}
		err = trello.UploadTrelloAttachment(cardId, downloadedFileName, &a, trelloKey, httpClient)
		if err != nil {
			log.Printf("ERROR Could not upload %s: %s", a.Filename, err)
			return err
		}
	}
	return nil
}
