package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

/*
writeDodgyContent is a debugging function that writes the given byte buffer to a file
*/
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

func LoadIssues(hostname string, key *ScriptKey, startAt int, pageSize int, maybeQuery string, httpClient *http.Client) (*PagedIssues, error) {
	uri := fmt.Sprintf("https://%s/rest/api/3/search?startAt=%d&maxResults=%d&fields=*all&expand=names", hostname, startAt, pageSize)
	if maybeQuery != "" {
		uri += "&jql=" + url.QueryEscape(maybeQuery)
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(key.User, key.Key)
	response, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bodyContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var issues PagedIssues

	switch response.StatusCode {
	case 200:
		//writeDodgyContent(fmt.Sprintf("content-%d.json", startAt/pageSize), &bodyContent)
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

func AsyncLoadIssuesJQL(hostname string, key *ScriptKey, pageSize int, maybeQuery string) (chan Issue, chan error) {
	outCh := make(chan Issue, 50)
	errCh := make(chan error, 1)

	go func() {
		ctr := 0
		httpClient := &http.Client{}
		for {
			pageData, err := LoadIssues(hostname, key, ctr, pageSize, maybeQuery, httpClient)
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

func AsyncLoadAllIssues(hostname string, key *ScriptKey, pageSize int) (chan Issue, chan error) {
	return AsyncLoadIssuesJQL(hostname, key, pageSize, "issueType in (Bug,Task,Story,Subtask)")
}

func AsyncLoadAllEpics(hostname string, key *ScriptKey, pageSize int) (chan Issue, chan error) {
	return AsyncLoadIssuesJQL(hostname, key, pageSize, "issueType=epic")
}

func SyncLoadAllEpics(hostname string, key *ScriptKey, pageSize int) ([]Issue, error) {
	outputCh, errCh := AsyncLoadAllEpics(hostname, key, pageSize)
	result := make([]Issue, 0)

	for {
		select {
		case rec, moreContent := <-outputCh:
			result = append(result, rec)
			if !moreContent {
				return result, nil
			}
		case err := <-errCh:
			return nil, err
		}
	}
}
