package trello

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fredex42/mm-jira-migration/common"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

type CustomFieldCache map[string]common.TrelloCustomField

func LoadAllCustomFields(boardId string, apiKey *common.ScriptKey, httpClient *http.Client) (*CustomFieldCache, error) {
	uri := fmt.Sprintf("https://api.trello.com/1/boards/%s/customFields?key=%s&token=%s", boardId, apiKey.User, apiKey.Key)
	response, err := httpClient.Get(uri)
	if err != nil {
		return nil, err
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	var customFieldList []common.TrelloCustomField
	switch response.StatusCode {
	case 200:
		err = json.Unmarshal(responseContent, &customFieldList)
		if err != nil {
			log.Printf("ERROR LoadAllCustomFields invalid response was %s", string(responseContent))
			log.Printf("ERROR LoadAllCustomFields could not understand server response: %s", err)
			return nil, err
		}
		output := make(CustomFieldCache, len(customFieldList))
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
func CreateCustomField(boardId string, name string, fieldType common.CustomFieldType, displayCardFront bool, options []string, apiKey *common.ScriptKey, httpClient *http.Client) (*common.TrelloCustomField, error) {
	uri := fmt.Sprintf("https://api.trello.com/1/customFields?key=%s&token=%s", apiKey.User, apiKey.Key)

	var optionsArg *string //defaults to 'nil'
	if fieldType == common.Checkbox {
		tempString := strings.Join(options, ",")
		optionsArg = &tempString
	}

	req := common.NewTrelloCustomField{
		BoardId:          boardId,
		ModelType:        "board",
		Name:             name,
		Type:             fieldType,
		Options:          optionsArg,
		Position:         common.TrelloPositionBottom,
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
		var out common.TrelloCustomField
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

func UpdateCustomField(definition *common.TrelloCustomField, apiKey *common.ScriptKey, httpClient *http.Client) (*common.TrelloCustomField, error) {
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
		var out common.TrelloCustomField
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

func AddCustomFieldOption(customFieldId string, definition *common.TrelloCustomFieldOption, apiKey *common.ScriptKey, httpClient *http.Client) error {
	uri := fmt.Sprintf("https://api.trello.com/1/customFields/%s/options?key=%s&token=%s", customFieldId, apiKey.User, apiKey.Key)

	bodyContent, err := json.Marshal(definition)
	if err != nil {
		return err
	}
	bodyContentReader := bytes.NewReader(bodyContent)
	response, err := httpClient.Post(uri, "application/json", bodyContentReader)
	if err != nil {
		return err
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("INFO Successfully added option to %s", customFieldId)
		return nil
	default:
		log.Printf("ERROR AddCustomFieldOption server response was %s", string(responseContent))
		msg := fmt.Sprintf("ERROR AddCustomFieldOption update custom field options on %s: Server error %d", customFieldId, response.StatusCode)
		return errors.New(msg)
	}
}

func RemoveCustomFieldOption(customFieldId string, optionId string, apiKey *common.ScriptKey, httpClient *http.Client) error {
	uri := fmt.Sprintf("https://api.trello.com/1/customFields/%s/options/%s?key=%s&token=%s", customFieldId, optionId, apiKey.User, apiKey.Key)

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return nil
	}
	response, err := httpClient.Do(req)
	if err != nil {
		return nil
	}
	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("INFO Successfully deleted option %s", optionId)
		return nil
	default:
		log.Printf("ERROR RemoveCustomFieldOption server response was %s", string(responseContent))
		msg := fmt.Sprintf("ERROR RemoveCustomFieldOption update custom field options on %s: Server error %d", customFieldId, response.StatusCode)
		return errors.New(msg)
	}
}

func SetupEpicsField(boardId string, customFieldName string, epicsList *[]common.Issue, apiKey *common.ScriptKey) (*common.TrelloCustomField, error) {
	httpClient := &http.Client{}
	existingCustomFields, err := LoadAllCustomFields(boardId, apiKey, httpClient)
	if err != nil {
		log.Printf("ERROR SetupEpicsField could not load existing fields: %s", err)
		return nil, err
	}

	var existingField *common.TrelloCustomField
	field, haveExistingField := (*existingCustomFields)[customFieldName]

	if haveExistingField {
		log.Printf("INFO SetupEpicsField Found existing field with name '%s'", customFieldName)
		existingField = &field
	} else { //we don't have an existing custom field
		log.Printf("INFO SetupEpicsField No existing field with name '%s', creating a new one...", customFieldName)
		opts := make([]string, len(*epicsList))
		//for i, e := range *epicsList {
		//	if e.Fields.EpicName != nil {
		//		log.Printf("INFO SetupEpicsField found epic name '%s'", *e.Fields.EpicName)
		//		opts[i] = *e.Fields.EpicName
		//	} else {
		//		log.Printf("WARN SetupEpicsField returned epic issue '%s' has no epic title", e.Fields.Summary)
		//	}
		//}
		existingField, err = CreateCustomField(boardId, customFieldName, common.List, true, opts, apiKey, httpClient)
		if err != nil {
			log.Printf("ERROR SetupEpicsField Unable to create field '%s': %s", customFieldName, err)
			return nil, err
		}
	}

	log.Printf("epics list length %d", len(*epicsList))

	sortedEpics := *epicsList
	sort.SliceStable(sortedEpics, func(i, j int) bool {
		firstName := sortedEpics[i].Fields.EpicName
		secondName := sortedEpics[j].Fields.EpicName
		if firstName == nil || secondName == nil {
			return false
		}
		switch strings.Compare(*firstName, *secondName) {
		case -1:
			return true
		default:
			return false
		}
	})

	log.Printf("sorted list length %d", len(sortedEpics))
	for i, e := range sortedEpics {
		if e.Fields.EpicName != nil {
			log.Printf("INFO SetupEpicsField found epic name '%s'", *e.Fields.EpicName)
			newOption := &common.TrelloCustomFieldOption{
				Id:            "",
				CustomFieldId: existingField.Id,
				Value: common.TrelloCustomFieldOptionValue{
					Text: *e.Fields.EpicName,
				},
				Colour: e.Fields.TranslateEpicColour(),
				Pos:    int64(i) * 10,
			}
			err = AddCustomFieldOption(existingField.Id, newOption, apiKey, httpClient)
			if err != nil {
				log.Printf("ERROR Unable to add custom field option: %s", err)
			}
		}
	}

	return existingField, nil
}
