package common

func (content *JiraContent) toTextBlock() string {
	accumulator := ""

	for _, block := range content.Content {
		for _, line := range block.Content {
			accumulator += line.Text + "\n"
		}
		accumulator += "\n"
	}
	return accumulator
}

/*
Create a new trello card request from the given issue.
Note that this does NOT bring over any attachments - that must be done seperately
*/
func (issue *Issue) toTrelloCard(inList string, putToTop bool) *NewTrelloCard {
	newPos := TrelloPositionBottom
	if putToTop {
		newPos = TrelloPositionTop
	}

	var dueComplete *bool
	if issue.Fields.Status.Name == "Done" {
		dueComplete = BoolPtr(true)
	}

	//PRIORITY should be a "custom field" in Trello
	//EPIC can map to a "component" "custom field" in Trello
	//what about comments?
	return &NewTrelloCard{
		ListId:      inList,
		Name:        issue.Fields.Summary,
		Description: (&(issue.Fields.Description)).toTextBlock(),
		Position:    newPos,
		DueDate:     issue.Fields.DueDate,
		Start:       nil,
		DueComplete: dueComplete,
		Members:     nil, //we need to cross-reference jira users to board members
		LabelIDs:    nil, //we need to merge in information about the epic, and any other labels
	}
}
