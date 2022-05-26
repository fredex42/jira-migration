package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fredex42/mm-jira-migration/common"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func writeDodgyContent(toFileName string, buffer *[]byte) error {
	f, err := os.OpenFile(toFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bytes.NewReader(*buffer)

	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}
	return nil
}

func LoadIssues(hostname string, key *common.ScriptKey, startAt int, pageSize int) (*common.PagedIssues, error) {
	uri := fmt.Sprintf("https://%s/rest/api/3/search?startAt=%d&maxResults=%d&fields=*all&expand=names", hostname, startAt, pageSize)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(key.User, key.Key)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bodyContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var issues common.PagedIssues

	switch response.StatusCode {
	case 200:
		err = json.Unmarshal(bodyContent, &issues)
		if err != nil {
			log.Printf("Unmarshalling error. Invalid content is being written to 'dodgy.json' ")
			writeDodgyContent("dodgy.json", &bodyContent)
			return nil, err
		}
		return &issues, nil
	default:
		log.Printf("Server returned %d. Body content was: ", response.StatusCode)
		log.Print(string(bodyContent))
		return nil, errors.New("server error")
	}
}

func AsyncLoadAllIssues(hostname string, key *common.ScriptKey, pageSize int) (chan common.Issue, chan error) {
	outCh := make(chan common.Issue, 50)
	errCh := make(chan error, 1)

	go func() {
		ctr := 0
		for {
			pageData, err := LoadIssues(hostname, key, ctr, pageSize)
			if err != nil {
				log.Printf("ERROR Can't load issues page: ")
				errCh <- err
				return
			}
			for _, i := range pageData.Issues {
				outCh <- i
			}
			ctr += len(pageData.Issues)
			if int64(ctr) >= pageData.Total {
				log.Printf("INFO Iterated a total of %d issues, completed", ctr)
				close(outCh)
				close(errCh)
				return
			}
		}
	}()

	return outCh, errCh
}

func main() {
	keyPath := flag.String("key", "scriptkey.yaml", "Path to a file containing a Jira API key")
	hostname := flag.String("host", "", "Virtual Jira host to query")
	pageSize := flag.Int("pagesize", 50, "number of issues to fetch in one page")
	flag.Parse()

	key, err := common.LoadScriptKey(keyPath)
	if err != nil {
		log.Fatalf("Could not open scripting key '%s': %s", *keyPath, err)
	}

	//content, err := LoadIssues(*hostname, key)
	//if err != nil {
	//	log.Fatal("Could not load content: ", err)
	//}
	//spew.Dump(content)

	contentCh, errCh := AsyncLoadAllIssues(*hostname, key, *pageSize)
	for {
		select {
		case err := <-errCh:
			log.Printf("ERROR: %s", err)
			os.Exit(1)
		case rec, moreContent := <-contentCh:
			log.Printf("INFO: Got record %v", rec)
			if !moreContent {
				return
			}
		}
	}
}
