package common

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
