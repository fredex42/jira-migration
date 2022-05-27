package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
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

func LoadAllCustomFields(boardId string, apiKey *ScriptKey, httpClient *http.Client) (*map[string]TrelloCustomField, error) {
	uri := fmt.Sprintf("https://api.trello.com/1/boards/%s/customFields?key=%s&token=%s", boardId, apiKey.User, apiKey.Key)
	response, err := httpClient.Get(uri)
	if err != nil {
		return nil, err
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	var customFieldList []TrelloCustomField
	switch response.StatusCode {
	case 200:
		err = json.Unmarshal(responseContent, &customFieldList)
		if err != nil {
			log.Printf("ERROR LoadAllCustomFields invalid response was %s", string(responseContent))
			log.Printf("ERROR LoadAllCustomFields could not understand server response: %s", err)
			return nil, err
		}
		output := make(map[string]TrelloCustomField, len(customFieldList))
		for _, f := range customFieldList {
			output[f.Name] = f
		}
		return &output, nil
	default:
		log.Printf("ERROR LoadAllCustomFields server response was %s", string(responseContent))
		msg := fmt.Sprintf("ERROR LoadAllCustomFields load custom fields on %s: Server error %d", boardId, response.StatusCode)
		return nil, errors.New(msg)
	}
}

/*
CreateCustomField creates a custom field on the given board, with the given parameters
*/
func CreateCustomField(boardId string, name string, fieldType CustomFieldType, displayCardFront bool, options []string, apiKey *ScriptKey, httpClient *http.Client) (*TrelloCustomField, error) {
	uri := fmt.Sprintf("https://api.trello.com/1/customFields?key=%s&token=%s", apiKey.User, apiKey.Key)

	var optionsArg *string //defaults to 'nil'
	if fieldType == Checkbox {
		tempString := strings.Join(options, ",")
		optionsArg = &tempString
	}

	req := NewTrelloCustomField{
		BoardId:          boardId,
		ModelType:        "board",
		Name:             name,
		Type:             fieldType,
		Options:          optionsArg,
		Position:         TrelloPositionBottom,
		DisplayCardFront: displayCardFront,
	}
	bodyContent, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}
	contentReader := bytes.NewReader(bodyContent)

	response, err := httpClient.Post(uri, "application/json", contentReader)
	if err != nil {
		return nil, err
	}
	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("INFO Successfully created field '%s' on board '%s'", name, boardId)
		var out TrelloCustomField
		err = json.Unmarshal(responseContent, &out)
		if err != nil {
			log.Print("ERROR CreateCustomField Invalid response was: ", string(responseContent))
			log.Printf("ERROR CreateCustomField Unable to understand server response: %s", err)
			return nil, err
		}
		return &out, nil
	default:
		log.Printf("ERROR CreateCustomField server response was %s", string(responseContent))
		msg := fmt.Sprintf("ERROR CreateCustomField load custom fields on %s: Server error %d", boardId, response.StatusCode)
		return nil, errors.New(msg)
	}
}

func UpdateCustomField(definition *TrelloCustomField, apiKey *ScriptKey, httpClient *http.Client) (*TrelloCustomField, error) {
	uri := fmt.Sprintf("https://api.trello.com/1/customFields/%s?key=%s&token=%s", definition.Id, apiKey.User, apiKey.Key)

	bodyContent, err := json.Marshal(definition)
	if err != nil {
		return nil, err
	}
	println(string(bodyContent))
	bodyContentReader := bytes.NewReader(bodyContent)

	req, err := http.NewRequest("PUT", uri, bodyContentReader)
	if err != nil {
		return nil, err
	}
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	println(string(responseContent))
	switch response.StatusCode {
	case 200:
		log.Printf("INFO Successfully updated field '%s' on board '%s'", definition.Name, definition.BoardId)
		var out TrelloCustomField
		err = json.Unmarshal(responseContent, &out)
		if err != nil {
			log.Print("ERROR CreateCustomField Invalid response was: ", string(responseContent))
			log.Printf("ERROR CreateCustomField Unable to understand server response: %s", err)
			return nil, err
		}
		return &out, nil
	default:
		log.Printf("ERROR CreateCustomField server response was %s", string(responseContent))
		msg := fmt.Sprintf("ERROR CreateCustomField updatge custom fields on %s: Server error %d", definition.BoardId, response.StatusCode)
		return nil, errors.New(msg)
	}
}

func makeTrelloSuitableId() string {
	uid, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("ERROR Unable to generate UUID: %s", err))
	}
	return strings.ReplaceAll(uid.String(), "-", "")
}

func SetupEpicsField(boardId string, customFieldName string, epicsList *[]Issue, apiKey *ScriptKey) (*TrelloCustomField, error) {
	httpClient := &http.Client{}
	existingCustomFields, err := LoadAllCustomFields(boardId, apiKey, httpClient)
	if err != nil {
		log.Printf("ERROR SetupEpicsField could not load existing fields: %s", err)
		return nil, err
	}

	var existingField *TrelloCustomField
	field, haveExistingField := (*existingCustomFields)[customFieldName]

	if haveExistingField {
		log.Printf("INFO SetupEpicsField Found existing field with name '%s'", customFieldName)
		existingField = &field
		newOptions := make([]TrelloCustomFieldOption, len(*epicsList))
		for i, e := range *epicsList {
			var textValue string
			if e.Fields.EpicName != nil {
				textValue = *e.Fields.EpicName
			}
			newOptions[i] = TrelloCustomFieldOption{
				Id:            makeTrelloSuitableId(),
				CustomFieldId: existingField.Id,
				Value: TrelloCustomFieldOptionValue{
					Text: textValue,
				},
				Colour: e.Fields.TranslateEpicColour(),
				Pos:    int64(i) * 10,
			}
		}
		spew.Dump(newOptions)
		existingField.Options = &newOptions
		spew.Dump(existingField)
		existingField, err = UpdateCustomField(existingField, apiKey, httpClient)

	} else { //we don't have an existing custom field
		log.Printf("INFO SetupEpicsField No existing field with name '%s', creating a new one...", customFieldName)
		opts := make([]string, len(*epicsList))
		for i, e := range *epicsList {
			if e.Fields.EpicName != nil {
				log.Printf("INFO SetupEpicsField found epic name '%s'", *e.Fields.EpicName)
				opts[i] = *e.Fields.EpicName
			} else {
				log.Printf("WARN SetupEpicsField returned epic issue '%s' has no epic title", e.Fields.Summary)
			}
		}
		existingField, err = CreateCustomField(boardId, customFieldName, List, true, opts, apiKey, httpClient)
		if err != nil {
			log.Printf("ERROR SetupEpicsField Unable to create field '%s': %s", customFieldName, err)
			return nil, err
		}
	}

	return existingField, nil
}
