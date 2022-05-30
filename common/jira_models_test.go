package common

import (
	"encoding/json"
	"testing"
)

func TestCommentsContent(t *testing.T) {
	testData := `{
        "version": 1,
        "type": "doc",
        "content": [
          {
            "type": "paragraph",
            "content": [
              {
                "type": "text",
                "text": "For some reason, when I tried to open "
              },
              {
                "type": "text",
                "text": "https://pluto.gnm.int/pluto-core/project/23706",
                "marks": [
                  {
                    "type": "link",
                    "attrs": {
                      "href": "https://pluto.gnm.int/pluto-core/project/23706"
                    }
                  }
                ]
              },
              {
                "type": "text",
                "text": " the conversion process worked. I am not sure why it worked on the machine I tested on and not on other machines. I will leave this for the moment and concentrate on the other files that are still not opening or being converted."
              }
            ]
          }
        ]
      }`
	var content JiraContent
	err := json.Unmarshal([]byte(testData), &content)
	if err != nil {
		t.Errorf("Could not unmarshal test data: %s", err)
	}
}

func TestComment(t *testing.T) {
	testData := `{
      "self": "https://codemill.atlassian.net/rest/api/3/issue/21717/comment/24878",
      "id": "24878",
      "author": {
        "self": "https://codemill.atlassian.net/rest/api/3/user?accountId=557058%3A8cc31b55-00e7-48f3-af1a-6351be68af14",
        "accountId": "557058:8cc31b55-00e7-48f3-af1a-6351be68af14",
        "emailAddress": "david.allison@guardian.co.uk",
        "avatarUrls": {
          "48x48": "https://secure.gravatar.com/avatar/fc192a4be7a5c69a9e2ca28c5263bfce?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FDA-3.png",
          "24x24": "https://secure.gravatar.com/avatar/fc192a4be7a5c69a9e2ca28c5263bfce?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FDA-3.png",
          "16x16": "https://secure.gravatar.com/avatar/fc192a4be7a5c69a9e2ca28c5263bfce?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FDA-3.png",
          "32x32": "https://secure.gravatar.com/avatar/fc192a4be7a5c69a9e2ca28c5263bfce?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FDA-3.png"
        },
        "displayName": "David Allison",
        "active": true,
        "timeZone": "Europe/Stockholm",
        "accountType": "atlassian"
      },
      "body": {
        "version": 1,
        "type": "doc",
        "content": [
          {
            "type": "paragraph",
            "content": [
              {
                "type": "text",
                "text": "For some reason, when I tried to open "
              },
              {
                "type": "text",
                "text": "https://pluto.gnm.int/pluto-core/project/23706",
                "marks": [
                  {
                    "type": "link",
                    "attrs": {
                      "href": "https://pluto.gnm.int/pluto-core/project/23706"
                    }
                  }
                ]
              },
              {
                "type": "text",
                "text": " the conversion process worked. I am not sure why it worked on the machine I tested on and not on other machines. I will leave this for the moment and concentrate on the other files that are still not opening or being converted."
              }
            ]
          }
        ]
      },
      "updateAuthor": {
        "self": "https://codemill.atlassian.net/rest/api/3/user?accountId=557058%3A8cc31b55-00e7-48f3-af1a-6351be68af14",
        "accountId": "557058:8cc31b55-00e7-48f3-af1a-6351be68af14",
        "emailAddress": "david.allison@guardian.co.uk",
        "avatarUrls": {
          "48x48": "https://secure.gravatar.com/avatar/fc192a4be7a5c69a9e2ca28c5263bfce?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FDA-3.png",
          "24x24": "https://secure.gravatar.com/avatar/fc192a4be7a5c69a9e2ca28c5263bfce?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FDA-3.png",
          "16x16": "https://secure.gravatar.com/avatar/fc192a4be7a5c69a9e2ca28c5263bfce?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FDA-3.png",
          "32x32": "https://secure.gravatar.com/avatar/fc192a4be7a5c69a9e2ca28c5263bfce?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FDA-3.png"
        },
        "displayName": "David Allison",
        "active": true,
        "timeZone": "Europe/Stockholm",
        "accountType": "atlassian"
      },
      "created": "2022-05-18T11:23:43.391+0200",
      "updated": "2022-05-18T11:23:43.391+0200",
      "jsdPublic": true
    }`

	var content Comment
	err := json.Unmarshal([]byte(testData), &content)
	if err != nil {
		t.Errorf("Could not unmarshal test data: %s", err)
	}
}
