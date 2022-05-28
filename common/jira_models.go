package common

import (
	"log"
	"regexp"
)

type PagedIssues struct {
	Expand     string  `json:"expand"`
	StartAt    int64   `json:"startAt"`
	MaxResults int64   `json:"maxResults"`
	Total      int64   `json:"total"`
	Issues     []Issue `json:"issues"`
}

type Issue struct {
	Expand string      `json:"expand"`
	Id     string      `json:"id"`
	Self   string      `json:"self"` //self-uri
	Key    string      `json:"key"`  //what it's referred to as in the UI
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Parent       *Issue        `json:"parent"` //reference to parent issue, if this is a subtask
	Priority     IssuePriority `json:"priority"`
	Labels       []string      `json:"labels"` //TBC schema
	TimeEstimate *int64        `json:"timeestimate"`
	Status       IssueStatus   `json:"status"`
	Creator      JiraUser      `json:"creator"`
	Subtasks     []Issue       `json:"subTasks"`
	Reporter     JiraUser      `json:"reporter"`
	IssueType    IssueType     `json:"issuetype"`
	Summary      string        `json:"summary"`
	Description  JiraContent   `json:"description"`
	Attachment   []Attachment  `json:"attachment"`
	DueDate      *string       `json:"duedate"`
	EpicLink     *string       `json:"customfield_10014"` //catchy name, huh? the id is unique to our jira *sigh*
	EpicName     *string       `json:"customfield_10011"` //only set on epics
	EpicColour   *string       `json:"customfield_10013"` //only set on epics. Use the decoding function to get a "sensible" colour name
	SprintLink   *[]SprintLink `json:"customfield_10020"`
}

/*
   {
     "id": 155,
     "name": "GNM Backlog (To be continued..",
     "state": "future",
     "boardId": 31,
     "goal": ""
   }
*/
type SprintLink struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	State   *string `json:"state"`
	BoardId int64   `json:"boardId"`
	Goal    *string `json:"goal"`
}

/*
TranslateEpicColour uses information from https://jira.atlassian.com/browse/JRACLOUD-59765 to translate an ID
in the form ghx-label-nnn to a 'normal' colour name and returns it as a string.
String can be empty if EpicColour is not set.
*/
func (i IssueFields) TranslateEpicColour() string {
	/*
		pink , yellow , lime , blue , black , orange , red , purple , sky , green
	*/

	if i.EpicColour == nil {
		return ""
	} else {
		xtractor := regexp.MustCompile("ghx-label-(\\d+)")
		parts := xtractor.FindAllStringSubmatch(*i.EpicColour, -1)
		if parts == nil {
			log.Printf("ERROR Can't translate epic colour '%s' as the format is not ghx-label-nnn", *i.EpicColour)
			return ""
		} else {
			switch parts[0][1] { //the extracted number
			case "1":
				return "black"
			case "2":
				return "yellow"
			case "14":
				return "orange"
			case "8":
				return "purple"
			case "4":
				return "blue"
			case "5":
				return "lime"
			case "13":
				return "green"
			case "12":
				return "black"
			case "3":
				return "yellow"
			case "9":
				return "pink"
			case "7":
				return "purple"
			case "10":
				return "sky"
			case "11":
				return "sky"
			case "6":
				return "lime"
			default:
				return ""
			}
		}
	}
}

type Attachment struct {
	Id       string   `json:"id"`
	Filename string   `json:"filename"`
	Author   JiraUser `json:"author"`
	Created  string   `json:"created"` //ISO timestamp with +nnnn zone
	Size     int64    `json:"size"`
	MimeType string   `json:"mimeType"`
	Content  string   `json:"content"` //URL of the attachment
}

type IssuePriority struct {
	Self    string `json:"self"`
	IconUrl string `json:"iconUrl"`
	Name    string `json:"name"`
	Id      string `json:"id"`
}

type IssueStatus struct {
	Self        string `json:"self"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Id          string `json:"id"`
}

type JiraUser struct {
	Self         string `json:"self"`
	AccountId    string `json:"accountId"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
}

type IssueType struct {
	Self           string `json:"self"`
	Id             string `json:"id"`
	Description    string `json:"description"`
	IconUrl        string `json:"iconUrl"`
	Name           string `json:"name"`
	SubTask        bool   `json:"subtask"`
	AvatarId       int64  `json:"avatarId"`
	HierarchyLevel int64  `json:"hierarchyLevel"`
}

type JiraContent struct {
	Version int32              `json:"version"`
	Type    string             `json:"type"`
	Content []JiraContentBlock `json:"content"`
}

type JiraContentBlock struct {
	Type    string            `json:"type"`
	Content []JiraContentLine `json:"content"`
}

type JiraContentLine struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
