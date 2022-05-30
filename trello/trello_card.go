package trello

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fredex42/mm-jira-migration/common"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func PutTrelloCard(definition *common.NewTrelloCard, apiKey *common.ScriptKey, httpClient *http.Client) (*common.TrelloCard, error) {
	uri := fmt.Sprintf("https://api.trello.com/1/cards?key=%s&token=%s", apiKey.User, apiKey.Key)
	uri += "&" + definition.ToQueryParams()

	log.Printf("URI is %s", uri)
	response, err := httpClient.Post(uri, "", nil)
	if err != nil {
		return nil, err
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == 200 {
		var out common.TrelloCard
		err = json.Unmarshal(responseContent, &out)
		if err != nil {
			log.Printf("ERROR Invalid content was %s", string(responseContent))
			return nil, errors.New("could not understand server response")
		}
		log.Printf("INFO PutTrelloCard created new card '%s' at '%s'", out.Name, out.ShortUrl)
		return &out, nil
	} else {
		log.Print("ERROR Server response was ", string(responseContent))
		msg := fmt.Sprintf("server returned %d when trying to create a card", response.StatusCode)
		return nil, errors.New(msg)
	}
}

/*
SetCustomFieldValue sets the value for a "list" type customfield on a card.

- cardId ID of the card to set
- fieldId ID of the customfield to set (this is NOT the name. Use trello.LoadAllCustomFields to get a CustomFieldCache to look up this value.
- value ID of the value to set. This must be a pre-existing value or the server returns 400. Look it up from the `options` list of a TrelloCustomField instance
- trelloKey common.ScriptKey pointer giving API credentials
- httpClient http client instance to use. This enables connection re-use
*/
func SetCustomFieldValue(cardId string, fieldId string, value string, trelloKey *common.ScriptKey, httpClient *http.Client) error {
	uri := fmt.Sprintf("https://api.trello.com/1/cards/%s/customField/%s/item?key=%s&token=%s&idValue=%s", cardId, fieldId, trelloKey.User, trelloKey.Key, url.QueryEscape(value))
	req, err := http.NewRequest("PUT", uri, nil)
	if err != nil {
		return err
	}

	return internalSetCustomField(req, httpClient)
}

func SetCustomFieldText(cardId string, fieldId string, value string, trelloKey *common.ScriptKey, httpClient *http.Client) error {
	uri := fmt.Sprintf("https://api.trello.com/1/cards/%s/customField/%s/item?key=%s&token=%s", cardId, fieldId, trelloKey.User, trelloKey.Key)

	contentDict := map[string]interface{}{
		"value": map[string]string{
			"text": value,
		},
	}

	contentBody, err := json.Marshal(&contentDict)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(contentBody)
	req, err := http.NewRequest("PUT", uri, reader)
	req.Header.Add("Content-Type", "application/json")
	return internalSetCustomField(req, httpClient)
}

func internalSetCustomField(req *http.Request, httpClient *http.Client) error {
	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	log.Printf("INFO SetupLinkField server returned %d %s", response.StatusCode, string(responseContent))

	if response.StatusCode == 200 {
		return nil
	} else {
		msg := fmt.Sprintf("server returned %d", response.StatusCode)
		return errors.New(msg)
	}
}
func AddComment(cardId string, content string, trelloKey *common.ScriptKey, httpClient *http.Client) error {
	uri := fmt.Sprintf("https://api.trello.com/1/cards/%s/actions/comments?key=%s&token=%s&text=%s", cardId, trelloKey.User, trelloKey.Key, url.QueryEscape(content))

	response, err := httpClient.Post(uri, "", nil)
	if err != nil {
		return err
	}
	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	log.Printf("INFO AddComment server returned %d", response.StatusCode)
	if response.StatusCode == 200 {
		return nil
	} else {
		log.Printf("ERROR AddComment server said %s", string(responseContent))
		return errors.New(fmt.Sprintf("server returned %d", response.StatusCode))
	}
}
