package common

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func StringPtr(str string) *string {
	return &str
}

func BoolPtr(b bool) *bool {
	return &b
}

//https://developer.atlassian.com/cloud/trello/rest/api-group-cards/#api-cards-post
//Note that all optional or nullable fields are built as pointers
type TrelloPosition string

const (
	//A number value is also acceptable

	//TrelloPositionTop tells the API to put this item at the top of the list
	TrelloPositionTop TrelloPosition = "top"
	//TrelloPositionBottom tells the API to put this item at the bottom of the list
	TrelloPositionBottom TrelloPosition = "bottom"
)

/*
NewTrelloCard is the data that the API expects when creating a new card
*/
type NewTrelloCard struct {
	ListId string `json:"idList"` //ID of the list to put it in. ^[0-9a-fA-F]{24}$

	Name        string         `json:"name"`
	Description string         `json:"desc"`
	Position    TrelloPosition `json:"pos"` //The position of the new card. 'top', 'bottom', or a positive float
	DueDate     *string        `json:"due"` //date format TBC
	Start       *string        `json:"start"`
	DueComplete *bool          `json:"dueComplete"`
	Members     []string       `json:"idMembers"`
	LabelIDs    []string       `json:"idLabels"`
}

/*
ToQueryParams converts a NewTrelloCard struct into a set of query parameters
*/
func (c *NewTrelloCard) ToQueryParams() string {
	params := []string{
		fmt.Sprintf("idList=%s", url.QueryEscape(c.ListId)),
		fmt.Sprintf("name=%s", url.QueryEscape(c.Name)),
		fmt.Sprintf("desc=%s", url.QueryEscape(c.Description)),
		fmt.Sprintf("pos=%s", url.QueryEscape(string(c.Position))),
	}
	if c.DueDate != nil {
		params = append(params, fmt.Sprintf("due=%s", url.QueryEscape(*c.DueDate)))
	}
	if c.DueComplete != nil {
		params = append(params, fmt.Sprintf("due=%t", *c.DueComplete))
	}
	return strings.Join(params, "&")
}

/*
TrelloCard is the data model for a card on Trello - see https://developer.atlassian.com/cloud/trello/rest/api-group-cards/#api-cards-post
*/
type TrelloCard struct {
	Id               string        `json:"id"`
	Address          *string       `json:"address"`
	CheckItemStates  []string      `json:"checkItemStates"`
	Closed           bool          `json:"closed"`
	Coordinates      *string       `json:"coordinates"`
	CreationMethod   *string       `json:"creationmethod"`
	DateLastActivity string        `json:"dateLastActivity"`
	Description      string        `json:"desc"`
	Due              *string       `json:"due"`
	DueReminder      *string       `json:"dueReminder"`
	Email            string        `json:"email"`
	LabelIDs         []interface{} `json:"idLabels"` //the interface could be a string or an instance of Label
	ListId           string        `json:"idList"`
	Members          []string      `json:"idMembers"`
	ShortId          int64         `json:"idShort"`
	Labels           []string      `json:"labels"`
	Name             string        `json:"name"`
	Position         float64       `json:"pos"`
	ShortLink        string        `json:"shortLink"`
	ShortUrl         string        `json:"shortUrl"`
	Subscribed       bool          `json:"subscribed"`
	URL              string        `json:"url"`
}

type TrelloLabel struct {
	Id          string  `json:"id"` //ID if this label
	BoardId     string  `json:"idBoard"`
	Name        string  `json:"name"`
	MaybeColour *string `json:"color"`
}

type CustomFieldType string

const (
	Checkbox CustomFieldType = "checkbox"
	List                     = "list"
	Number                   = "number"
	Text                     = "text"
	Date                     = "date"
)

//NewTrelloCustomField describes the JSON document used for creating custom fields
type NewTrelloCustomField struct {
	BoardId          string          `json:"idModel"`
	ModelType        string          `json:"modelType"` //must be "board"
	Name             string          `json:"name"`      //name of the custom field
	Type             CustomFieldType `json:"type"`
	Options          *string         `json:"options"` //options for checkbox type. Assume `name` only, or it might be a documentation error
	Position         TrelloPosition  `json:"pos"`
	DisplayCardFront bool            `json:"display_cardFront"`
}

type TrelloCustomFieldDisplay struct {
	CardFront bool `json:"cardFront"`
}

type TrelloCustomFieldOptionValue struct {
	Text string `json:"text"`
}

type TrelloCustomFieldOption struct {
	Id            string                       `json:"id"`
	CustomFieldId string                       `json:"idCustomField"`
	Value         TrelloCustomFieldOptionValue `json:"value"`
	Colour        string                       `json:"color"`
	Pos           int64                        `json:"pos"`
}

type TrelloCustomField struct {
	Id               string                     `json:"id"`
	BoardId          string                     `json:"idModel"`
	ModelType        string                     `json:"modelType"` //must be "board"
	FieldGroup       string                     `json:"fieldGroup"`
	Display          TrelloCustomFieldDisplay   `json:"display"`
	Name             string                     `json:"name"`
	Pos              int64                      `json:"pos"`
	Options          *[]TrelloCustomFieldOption `json:"options"`
	Type             CustomFieldType            `json:"type"`
	IsSuggestedField bool                       `json:"isSuggestedField"`
}

func (f *TrelloCustomField) FindInCustomField(nameToFind string) (*TrelloCustomFieldOption, error) {
	for _, opt := range *f.Options {
		if opt.Value.Text == nameToFind {
			copied := opt
			return &copied, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("there is no custom field value matching '%s'", nameToFind))
}

type TrelloList struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Closed     bool              `json:"closed"`
	Pos        int64             `json:"pos"`
	SoftLimit  *string           `json:"softLimit"`
	BoardId    string            `json:"idBoard"`
	Subscribed bool              `json:"subscribed"`
	Limits     *ListLimitsObject `json:"limits"`
}

type ListLimitsObject struct {
	Attachments LimitsPerBoard `json:"attachments"`
}

type LimitsPerBoard struct {
	PerBoard LimitsObject `json:"perBoard"`
}

type LimitStatus string

const (
	Ok      LimitStatus = "ok"
	Warning             = "warning"
)

type LimitsObject struct {
	Status    LimitStatus `json:"status"`
	DisableAt *int64      `json:"disableAt"`
	WarnAt    *int64      `json:"warnAt"`
}
