package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

/*
CreateLabel Calls the Trello REST API to create a label

boardId: the ID of the board to target
name: name of the new label
maybeColour: either a colour name or the string "null"
apiKey: ScriptKey struct with the API key to use
*/
func CreateLabel(boardId string, name string, maybeColour string, apiKey *ScriptKey) (*TrelloLabel, error) {
	uri := fmt.Sprintf("https://api.trello.com/1/boards/%s/labels?name=%s&color=%s&key=%s&token=%s",
		boardId,
		url.QueryEscape(name),
		url.QueryEscape(maybeColour),
		apiKey.User,
		apiKey.Key)

	response, err := http.Post(uri, "", nil)
	if err != nil {
		return nil, err
	}
	responseContent, err := ioutil.ReadAll(response.Body)
	var label TrelloLabel

	switch response.StatusCode {
	case 200:
		err = json.Unmarshal(responseContent, &label)
		if err != nil {
			log.Printf("ERROR CreateLabel invalid response was %s", string(responseContent))
			log.Printf("ERROR CreateLabel could not understand server response: %s", err)
			return nil, err
		}
		return &label, nil
	default:
		log.Printf("ERROR CreateLabel server response was %s", string(responseContent))
		msg := fmt.Sprintf("ERROR CreateLabel could not create a label for %s on %s: Server error %d", name, boardId, response.StatusCode)
		return nil, errors.New(msg)
	}
}
