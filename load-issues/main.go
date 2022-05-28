package main

import (
	"flag"
	"github.com/davecgh/go-spew/spew"
	"github.com/fredex42/mm-jira-migration/common"
	"log"
	"os"
)

func main() {
	keyPath := flag.String("key", "scriptkey.yaml", "Path to a file containing a Jira API key")
	hostname := flag.String("host", "", "Virtual Jira host to query")
	pageSize := flag.Int("pagesize", 50, "number of issues to fetch in one page")
	flag.Parse()

	key, err := common.LoadScriptKey(keyPath)
	if err != nil {
		log.Fatalf("Could not open scripting key '%s': %s", *keyPath, err)
	}

	contentCh, errCh := common.AsyncLoadAllIssues(*hostname, key, *pageSize)
	for {
		select {
		case err := <-errCh:
			log.Printf("ERROR: %s", err)
			os.Exit(1)
		case rec, moreContent := <-contentCh:
			//if rec.Fields.EpicLink != nil {
			//	spew.Dump(rec)
			//	//log.Printf("INFO: Got record %v", rec)
			//}
			if rec.Fields.Status.Name == "Done" { //don't bother importing over stuff marked as 'Done'
				continue
			}
			if rec.Fields.SprintLink != nil {
				spew.Dump(rec)
			}
			if !moreContent {
				return
			}
		}
	}
}
