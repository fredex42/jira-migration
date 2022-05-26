package common

func StringPtr(str string) *string {
	return &str
}

func BoolPtr(b bool) *bool {
	return &b
}

//https://developer.atlassian.com/cloud/trello/rest/api-group-cards/#api-cards-post
//Note that all optional or nullable fields are built as pointers

type NewTrelloCard struct {
	ListId string `json:"idList"` //ID of the list to put it in. ^[0-9a-fA-F]{24}$

	Name        string   `json:"name"`
	Description string   `json:"desc"`
	Position    string   `json:"pos"` //The position of the new card. 'top', 'bottom', or a positive float
	DueDate     *string  `json:"due"` //date format TBC
	Start       *string  `json:"start"`
	DueComplete *bool    `json:"dueComplete"`
	Members     []string `json:"idMembers"`
	LabelIDs    []string `json:"idLabels"`
}

type TrelloLabel struct {
	Id          string  `json:"id"` //ID if this label
	BoardId     string  `json:"idBoard"`
	Name        string  `json:"name"`
	MaybeColour *string `json:"color"`
}
