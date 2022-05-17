package common

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
	TimeEstimate *string       `json:"timeestimate"`
	Status       IssueStatus   `json:"status"`
	Creator      JiraUser      `json:"creator"`
	Subtasks     []Issue       `json:"subTasks"`
	Reporter     JiraUser      `json:"reporter"`
	IssueType    IssueType     `json:"issuetype"`
	Summary      string        `json:"summary"`
	Description  JiraContent   `json:"description"`
	Attachment   []string      `json:"attachment"` //TBC schema
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
