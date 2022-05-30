package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/fredex42/mm-jira-migration/common"
	"github.com/fredex42/mm-jira-migration/trello"
	"log"
	"net/http"
	"os"
	"time"
)

/*
MakeEpicLink sets the custom field on a created trello card to the epic's value.
Assumes that recPtr.Fields.EpicLink != nil, will abort if this is not the case.
*/
func MakeEpicLink(recPtr *common.Issue, cardId string, epics *EpicsCache, epicLinkField *common.TrelloCustomField, trelloKey *common.ScriptKey, httpClient *http.Client) error {
	log.Printf("INFO Issue '%s' has a link to epic '%s'", recPtr.Fields.Summary, *recPtr.Fields.EpicLink)

	epicName, haveEpic := epics.KnownEpics[*(recPtr.Fields.EpicLink)]
	if !haveEpic {
		log.Printf("ERROR Could not find an epic for '%s'", *(recPtr.Fields.EpicLink))
		return errors.New("could not create epic link")
	}

	epicId, err := epicLinkField.FindInCustomField(epicName)
	if err != nil {
		log.Printf("ERROR Could not find an entry for epic '%s' %s: %s", *recPtr.Fields.EpicLink, epicName, err)
		return errors.New("could not create epic link")
	}

	err = trello.SetCustomFieldValue(cardId, epicLinkField.Id, epicId.Id, trelloKey, httpClient) //should use TrelloCustomFieldOptionValue as k-v i think. https://developer.atlassian.com/cloud/trello/rest/api-group-cards/#api-cards-idcard-customfield-idcustomfield-item-put
	if err != nil {
		log.Printf("ERROR Could not set up custom epics info field for '%s': %s", recPtr.Fields.Summary, err)
		return errors.New("could not create epic link")
	}
	return nil
}

func MigrateIssue(recPtr *common.Issue,
	hostname *string,
	defaultList *common.TrelloList,
	epicLinkField *common.TrelloCustomField,
	priorityField *common.TrelloCustomField,
	epics *EpicsCache,
	jiraIdField *common.TrelloCustomField,
	jira *common.ScriptKey,
	trelloKey *common.ScriptKey,
	httpClient *http.Client,
) error {
	//get a base trello card
	newCard := recPtr.ToTrelloCard(defaultList.Id, false)
	//write the card and get an ID
	createdCard, err := trello.PutTrelloCard(newCard, trelloKey, httpClient)
	if err != nil {
		log.Printf("ERROR Could not create a card for '%s': %s", recPtr.Fields.Summary, err)
		return errors.New("can't migrate issue")
	}

	err = trello.SetCustomFieldText(createdCard.Id, jiraIdField.Id, recPtr.Key, trelloKey, httpClient)
	if err != nil {
		log.Printf("ERROR Could not add jira key for '%s': %s", recPtr.Fields.Summary, err)
		return errors.New("can't migrate issue")
	}

	//if there are attachments, copy them over
	err = HandleAttachments(&recPtr.Fields.Attachment, createdCard.Id, hostname, jira, trelloKey, httpClient)
	if err != nil {
		log.Printf("ERROR Could not fix attachments for '%s': %s", recPtr.Fields.Summary, err)
		return errors.New("can't migrate issue")
	}
	//if there is an epic link, find the custom field value corresponding and set it
	if recPtr.Fields.EpicLink != nil {
		err = MakeEpicLink(recPtr, createdCard.Id, epics, epicLinkField, trelloKey, httpClient)
		if err != nil {
			return errors.New("can't migrate issue")
		}
	}

	//get the priority and set that on a custom field too
	fieldId, err := recPtr.Fields.Priority.ToTrelloLabel(priorityField.Options)
	if err != nil {
		log.Printf("ERROR Could not set up priority for '%s': '%s", recPtr.Fields.Summary, err)
		return errors.New("can't migrate issue")
	}
	err = trello.SetCustomFieldValue(createdCard.Id, priorityField.Id, fieldId, trelloKey, httpClient)

	if err != nil {
		log.Printf("ERROR Could not set up priority field for '%s': %s", recPtr.Fields.Summary, err)
		return errors.New("can't migrate issue")
	}

	//now need to migrate all other comments
	existingComments, err := common.LoadAllComments(*hostname, recPtr.Key, jira, 20, httpClient)
	if err != nil {
		log.Printf("ERROR Can't load comments for '%s': %s", recPtr.Fields.Summary, err)
		return errors.New("can't migrate issue")
	}

	for _, c := range *existingComments {
		createdTime, err := time.Parse(common.JiraTimeFormat, c.Created)
		var createdTimeString string

		if err != nil {
			log.Printf("ERROR Can't parse time %s: %s", c.Created, err)
			createdTimeString = c.Created
		} else {
			createdTimeString = createdTime.Format(time.RFC1123)
		}

		newComment := fmt.Sprintf("%s\n-----\nOriginally by %s on %s", c.Body.ToTextBlock(), c.Author.DisplayName, createdTimeString)
		err = trello.AddComment(createdCard.Id, newComment, trelloKey, httpClient)
		if err != nil {
			log.Printf("ERROR Could not add comment to card '%s': %s", createdCard.Id, err)
			return errors.New("can't migrate issue")
		}
	}

	//set a comment showing where this came from and when
	creationTime, parseErr := time.Parse(common.JiraTimeFormat, recPtr.Fields.Created)
	creationTimeString := recPtr.Fields.Created
	if parseErr != nil {
		log.Printf("WARNING Can't parse time '%s': %s", recPtr.Fields.Created, parseErr)
	} else {
		creationTimeString = creationTime.Format(time.RFC1123)
	}

	newComment := fmt.Sprintf(`This card was originally reported on %s by %s as issue %s`,
		creationTimeString,
		recPtr.Fields.Reporter.DisplayName,
		recPtr.Key,
	)
	err = trello.AddComment(createdCard.Id, newComment, trelloKey, httpClient)
	if err != nil {
		log.Printf("ERROR Could not add comment to card '%s': %s", createdCard.Id, err)
		return errors.New("can't migrate issue")
	}
	return nil
}

func main() {
	jiraKeyPath := flag.String("jira", "scriptkey.yaml", "Path to a file containing a Jira API key")
	trelloKeyPath := flag.String("trello", "trellokey.yaml", "Path to a file containing a Trello API key")
	hostname := flag.String("host", "", "Virtual Jira host to query")
	pageSize := flag.Int("pagesize", 50, "number of issues to fetch in one page")
	trelloBoard := flag.String("board", "", "Board ID to push data into")
	defaultList := flag.String("defaultlist", "", "Name of the list to push cards into by default")
	epicLinkFieldName := flag.String("epicfield", "Components", "Name of the custom field to hold epics information")
	jiraIdFieldName := flag.String("jira-id", "Jira Key", "Name of the custom field to hold the jira ID")
	flag.Parse()

	httpClient := http.DefaultClient

	jira, err := common.LoadScriptKey(jiraKeyPath)
	if err != nil {
		log.Fatalf("Could not open scripting key '%s': %s", *jiraKeyPath, err)
	}

	trelloKey, err := common.LoadScriptKey(trelloKeyPath)
	if err != nil {
		log.Fatalf("Could not open scripting key '%s': %s", *trelloKeyPath, err)
	}

	trelloListCache, err := trello.NewListCache(*trelloBoard, trelloKey, httpClient)
	if err != nil {
		log.Fatalf("Could not load lists from board '%s': %s", *trelloBoard, err)
	}
	log.Printf("INFO Found %d lists on board '%s' ", trelloListCache.Count(), *trelloBoard)

	customFieldCache, err := trello.LoadAllCustomFields(*trelloBoard, trelloKey, httpClient)
	if err != nil {
		log.Fatalf("Could not load custom fields from board '%s': %s", *trelloBoard, err)
	}
	log.Printf("INFO Found %d custom fields on board '%s'", len(*customFieldCache), *trelloBoard)

	defaultListId, haveList := trelloListCache.FindByName(*defaultList)
	if !haveList {
		log.Fatalf("There is no list '%s' on the board", *defaultList)
	}

	epicLinkField, haveEpicLinkField := (*customFieldCache)[*epicLinkFieldName]
	if !haveEpicLinkField {
		log.Fatalf("Could not find any custom field matching '%s' for epics information", *epicLinkFieldName)
	}

	jiraIdField, haveJiraIdField := (*customFieldCache)[*jiraIdFieldName]
	if !haveJiraIdField {
		log.Fatalf("Could not find any custom field matching '%s' for jira key information", *jiraIdFieldName)
	}

	priorityField, havePriorityField := (*customFieldCache)["Priority"]
	if !havePriorityField {
		log.Fatal("Could not find any custom field matching 'Priority' for priority information")
	}

	epics, err := NewEpicsCache(hostname, jira, *pageSize)
	if err != nil {
		log.Fatal("Unable to load epics information")
	}
	if len(epics.KnownEpics) == 0 {
		log.Fatal("Could not load in any epics, check the code")
	}

	contentCh, errCh := common.AsyncLoadAllIssues(*hostname, jira, *pageSize)
	ctr := 0

	for {
		select {
		case err := <-errCh:
			log.Printf("ERROR: %s", err)
			os.Exit(1)
		case rec, moreContent := <-contentCh:
			if rec.Fields.Status.Name == "Done" { //don't bother importing over stuff marked as 'Done'
				continue
			}

			err := MigrateIssue(&rec, hostname, &defaultListId, &epicLinkField, &priorityField, epics, &jiraIdField, jira, trelloKey, httpClient)
			if err != nil {
				log.Fatalf("ERROR processing '%s': %s", rec.Key, err)
			}

			ctr++

			if !moreContent {
				log.Printf("Job completed! Migrated %d issues over", ctr)
				return
			}
		}
	}
}
