package trello

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fredex42/mm-jira-migration/common"
	"io/ioutil"
	"log"
	"net/http"
)

type ListCache struct {
	BoardId string
	knownLists map[string]common.TrelloList
	listsById map[string]common.TrelloList
}

func GetListsForBoard(boardId string, apiKey *common.ScriptKey, httpClient *http.Client) ([]common.TrelloList, error) {
	uri := fmt.Sprintf("https://api.trello.com/1/boards/%s/lists?key=%s&token=%s", boardId, apiKey.User, apiKey.Key)

	response, err := httpClient.Get(uri)
	if err != nil {
		return nil, err
	}
	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode==200 {
		var out []common.TrelloList
		err = json.Unmarshal(responseContent, &out)
		if err != nil {
			return nil, err
		} else {
			return out, nil
		}
	} else {
		log.Printf("ERROR GetListsForBoard server response was %s", string(responseContent))
		msg := fmt.Sprintf("GetListsForBoard server returned %d", response.StatusCode)
		return nil, errors.New(msg)
	}
}

func NewListCache(boardId string, apiKey *common.ScriptKey, httpClient *http.Client) (*ListCache, error) {
	content, err := GetListsForBoard(boardId, apiKey, httpClient)
	if err != nil {
		return nil, err
	}

	cache := &ListCache{
		BoardId:    boardId,
		knownLists: make(map[string]common.TrelloList, len(content)),
		listsById: make(map[string]common.TrelloList, len(content)),
	}
	for _, l := range content {
		cache.knownLists[l.Name] = l
		cache.listsById[l.Id] = l
	}
	return cache, nil
}

func (c *ListCache) FindByName(listName string) (common.TrelloList, bool){
	content, haveResult := c.knownLists[listName]
	return content, haveResult
}

func (c *ListCache) FindById(listId string) (common.TrelloList, bool) {
	content, haveResult := c.listsById[listId]
	return content, haveResult
}

func (c *ListCache) Count() int {
	return len(c.knownLists)
}