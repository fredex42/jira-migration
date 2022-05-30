package trello

import (
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

func SetCustomFieldValue(cardId string, fieldId string, value string, trelloKey *common.ScriptKey, httpClient *http.Client) error {
	uri := fmt.Sprintf("https://api.trello.com/1/cards/%s/customField/%s/item?key=%s&token=%s&idValue=%s", cardId, fieldId, trelloKey.User, trelloKey.Key, url.QueryEscape(value))

	//contentDict := map[string]string{
	//	"idValue": value,
	//}
	//
	//bodyContent, err := json.Marshal(&contentDict)
	//if err != nil {
	//	return err
	//}
	//
	//log.Printf("SetupCustomFieldValue: body content is %s", string(bodyContent))
	//contentReader := bytes.NewReader(bodyContent)

	req, err := http.NewRequest("PUT", uri, nil)
	//req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

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
