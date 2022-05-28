package main

import (
	"flag"
	"github.com/fredex42/mm-jira-migration/common"
	"github.com/fredex42/mm-jira-migration/trello"
	"log"
)

func main() {
	jiraKeyPath := flag.String("jirakey", "scriptkey.yaml", "Path to a file containing a Jira API key")
	trelloKeyPath := flag.String("trellokey", "trellokey.yaml", "Path to a file containing a Trello API key")
	hostname := flag.String("host", "", "Virtual Jira host to query")
	pageSize := flag.Int("pagesize", 50, "number of issues to fetch in one page")
	boardId := flag.String("board", "", "Trello board to update")
	customFieldName := flag.String("field", "component", "Custom field to create or update with epic names")
	flag.Parse()

	jiraKey, err := common.LoadScriptKey(jiraKeyPath)
	if err != nil {
		log.Fatal("ERROR Could not load key from ", *jiraKeyPath, ": ", err)
	}
	trelloKey, err := common.LoadScriptKey(trelloKeyPath)
	if err != nil {
		log.Fatal("ERROR Could not load key from ", *trelloKeyPath, ": ", err)
	}

	epicsList, err := common.SyncLoadAllEpics(*hostname, jiraKey, *pageSize)
	if err != nil {
		log.Fatal("ERROR Could not load in epics: ", err)
	}

	fieldContent, err := trello.SetupEpicsField(*boardId, *customFieldName, &epicsList, trelloKey)
	if err != nil {
		log.Fatal("ERROR Could not upload content to Trello: ", err)
	}

	log.Printf("INFO Updated custom field '%s' on board '%s' with %d options", fieldContent.Name, fieldContent.BoardId, len(*fieldContent.Options))
}
