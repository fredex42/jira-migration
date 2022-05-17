package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/fredex42/mm-jira-migration/common"
	"io/ioutil"
	"log"
	"net/http"
)

func LoadIssues(hostname string, key *common.ScriptKey) (*common.PagedIssues, error) {
	uri := fmt.Sprintf("https://%s/rest/api/3/search", hostname)
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
			return nil, err
		}
		return &issues, nil
	default:
		log.Printf("Server returned %d. Body content was: ", response.StatusCode)
		log.Print(string(bodyContent))
		return nil, errors.New("server error")
	}
}

func main() {
	keyPath := flag.String("key", "scriptkey.yaml", "Path to a file containing a Jira API key")
	hostname := flag.String("host", "", "Virtual Jira host to query")
	flag.Parse()

	key, err := common.LoadScriptKey(keyPath)
	if err != nil {
		log.Fatalf("Could not open scripting key '%s': %s", *keyPath, err)
	}

	content, err := LoadIssues(*hostname, key)
	if err != nil {
		log.Fatal("Could not load content: ", err)
	}
	spew.Dump(content)
}
