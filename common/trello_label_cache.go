package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TrelloLabelCache struct {
	Labels map[string]TrelloLabel
}

/*
TrelloLabelCache.lookup returns the TrelloLabel with the given name or an empty trello label with false
*/
func (c *TrelloLabelCache) lookup(name string) (TrelloLabel, bool) {
	content, haveContent := c.Labels[name]
	return content, haveContent
}

/*
NewTrelloLabelCache initialises a new label cache object with the label contents of the given board
*/
func NewTrelloLabelCache(boardId string, key *ScriptKey) (*TrelloLabelCache, error) {
	uri := fmt.Sprintf("https://api.trello.com/1/boards/%s/labels?key=%s&token=%s", boardId, key.User, key.Key)

	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	contentBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var rawList []TrelloLabel
	err = json.Unmarshal(contentBody, &rawList)
	if err != nil {
		return nil, err
	}

	cache := TrelloLabelCache{
		Labels: make(map[string]TrelloLabel, len(rawList)),
	}

	for _, l := range rawList {
		cache.Labels[l.Name] = l
	}
	return &cache, nil
}
