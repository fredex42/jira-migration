package main

import (
	"flag"
	"github.com/fredex42/mm-jira-migration/common"
	"github.com/fredex42/mm-jira-migration/trello"
	"log"
	"net/http"
	"os"
)

func main() {
	jiraKeyPath := flag.String("jira", "scriptkey.yaml", "Path to a file containing a Jira API key")
	trelloKeyPath := flag.String("trello", "trellokey.yaml", "Path to a file containing a Trello API key")
	hostname := flag.String("host", "", "Virtual Jira host to query")
	pageSize := flag.Int("pagesize", 50, "number of issues to fetch in one page")
	trelloBoard := flag.String("board", "", "Board ID to push data into")
	defaultList := flag.String("defaultlist", "", "Name of the list to push cards into by default")
	epicLinkFieldName := flag.String("epicfield", "Components", "Name of the custom field to hold epics information")
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
			if rec.Fields.EpicLink == nil { //testing - only stuff with epic links
				continue
			}
			recPtr := &rec

			//get a base trello card
			newCard := recPtr.ToTrelloCard(defaultListId.Id, false)
			//write the card and get an ID
			createdCard, err := trello.PutTrelloCard(newCard, trelloKey, httpClient)
			if err != nil {
				log.Printf("ERROR Could not create a card for '%s': %s", recPtr.Fields.Summary, err)
				return //bail on error
			}

			//if there is an epic link, find the custom field value corresponding and set it
			if rec.Fields.EpicLink != nil {
				log.Printf("INFO Issue '%s' has a link to epic '%s'", recPtr.Fields.Summary, *rec.Fields.EpicLink)

				epicName, haveEpic := epics.KnownEpics[*rec.Fields.EpicLink]
				if !haveEpic {
					log.Printf("ERROR Could not find an epic for '%s'", *rec.Fields.EpicLink)
					return
				}

				epicId, err := epicLinkField.FindInCustomField(epicName)
				if err != nil {
					log.Printf("ERROR Could not find an entry for epic '%s' %s: %s", *rec.Fields.EpicLink, epicName, err)
					return
				}

				err = trello.SetCustomFieldValue(createdCard.Id, epicLinkField.Id, epicId.Id, trelloKey, httpClient) //should use TrelloCustomFieldOptionValue as k-v i think. https://developer.atlassian.com/cloud/trello/rest/api-group-cards/#api-cards-idcard-customfield-idcustomfield-item-put
				if err != nil {
					log.Printf("ERROR Could not set up custom epics info field for '%s': %s", recPtr.Fields.Summary, err)
					return
				}
			}

			//get the priority and set that on a custom field too
			fieldId, err := rec.Fields.Priority.ToTrelloLabel(priorityField.Options)
			if err != nil {
				log.Printf("ERROR Could not set up priority for '%s': '%s", recPtr.Fields.Summary, err)
				return
			}
			err = trello.SetCustomFieldValue(createdCard.Id, priorityField.Id, fieldId, trelloKey, httpClient)

			if err != nil {
				log.Printf("ERROR Could not set up priority field for '%s': %s", recPtr.Fields.Summary, err)
				return
			}

			ctr++
			if ctr > 3 {
				log.Printf("INFO Early exit for debugging")
				return
			}
			if !moreContent {
				return
			}
		}
	}
}
