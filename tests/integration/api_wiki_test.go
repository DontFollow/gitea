// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIGetWikiPage(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	username := "user2"
	session := loginUser(t, username)

	urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/Home", username, "repo1")

	req := NewRequest(t, "GET", urlStr)
	resp := session.MakeRequest(t, req, http.StatusOK)
	var page *api.WikiPage
	DecodeJSON(t, resp, &page)

	assert.Equal(t, &api.WikiPage{
		WikiPageMetaData: &api.WikiPageMetaData{
			Title:   "Home",
			HTMLURL: page.HTMLURL,
			SubURL:  "Home",
			LastCommit: &api.WikiCommit{
				ID: "2c54faec6c45d31c1abfaecdab471eac6633738a",
				Author: &api.CommitUser{
					Identity: api.Identity{
						Name:  "Ethan Koenig",
						Email: "ethantkoenig@gmail.com",
					},
					Date: "2017-11-27T04:31:18Z",
				},
				Committer: &api.CommitUser{
					Identity: api.Identity{
						Name:  "Ethan Koenig",
						Email: "ethantkoenig@gmail.com",
					},
					Date: "2017-11-27T04:31:18Z",
				},
				Message: "Add Home.md\n",
			},
		},
		ContentBase64: base64.RawStdEncoding.EncodeToString(
			[]byte("# Home page\n\nThis is the home page!\n"),
		),
		CommitCount: 1,
		Sidebar:     "",
		Footer:      "",
	}, page)
}

func TestAPIListWikiPages(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	username := "user2"
	session := loginUser(t, username)

	urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/pages", username, "repo1")

	req := NewRequest(t, "GET", urlStr)
	resp := session.MakeRequest(t, req, http.StatusOK)

	var meta []*api.WikiPageMetaData
	DecodeJSON(t, resp, &meta)

	dummymeta := []*api.WikiPageMetaData{
		{
			Title:   "Home",
			HTMLURL: meta[0].HTMLURL,
			SubURL:  "Home",
			LastCommit: &api.WikiCommit{
				ID: "2c54faec6c45d31c1abfaecdab471eac6633738a",
				Author: &api.CommitUser{
					Identity: api.Identity{
						Name:  "Ethan Koenig",
						Email: "ethantkoenig@gmail.com",
					},
					Date: "2017-11-27T04:31:18Z",
				},
				Committer: &api.CommitUser{
					Identity: api.Identity{
						Name:  "Ethan Koenig",
						Email: "ethantkoenig@gmail.com",
					},
					Date: "2017-11-27T04:31:18Z",
				},
				Message: "Add Home.md\n",
			},
		},
		{
			Title:   "Page With Image",
			HTMLURL: meta[1].HTMLURL,
			SubURL:  "Page-With-Image",
			LastCommit: &api.WikiCommit{
				ID: "0cf15c3f66ec8384480ed9c3cf87c9e97fbb0ec3",
				Author: &api.CommitUser{
					Identity: api.Identity{
						Name:  "Gabriel Silva Simões",
						Email: "simoes.sgabriel@gmail.com",
					},
					Date: "2019-01-25T01:41:55Z",
				},
				Committer: &api.CommitUser{
					Identity: api.Identity{
						Name:  "Gabriel Silva Simões",
						Email: "simoes.sgabriel@gmail.com",
					},
					Date: "2019-01-25T01:41:55Z",
				},
				Message: "Add jpeg.jpg and page with image\n",
			},
		},
		{
			Title:   "Page With Spaced Name",
			HTMLURL: meta[2].HTMLURL,
			SubURL:  "Page-With-Spaced-Name",
			LastCommit: &api.WikiCommit{
				ID: "c10d10b7e655b3dab1f53176db57c8219a5488d6",
				Author: &api.CommitUser{
					Identity: api.Identity{
						Name:  "Gabriel Silva Simões",
						Email: "simoes.sgabriel@gmail.com",
					},
					Date: "2019-01-25T01:39:51Z",
				},
				Committer: &api.CommitUser{
					Identity: api.Identity{
						Name:  "Gabriel Silva Simões",
						Email: "simoes.sgabriel@gmail.com",
					},
					Date: "2019-01-25T01:39:51Z",
				},
				Message: "Add page with spaced name\n",
			},
		},
		{
			Title:   "Unescaped File",
			HTMLURL: meta[3].HTMLURL,
			SubURL:  "Unescaped-File",
			LastCommit: &api.WikiCommit{
				ID: "0dca5bd9b5d7ef937710e056f575e86c0184ba85",
				Author: &api.CommitUser{
					Identity: api.Identity{
						Name:  "6543",
						Email: "6543@obermui.de",
					},
					Date: "2021-07-19T16:42:46Z",
				},
				Committer: &api.CommitUser{
					Identity: api.Identity{
						Name:  "6543",
						Email: "6543@obermui.de",
					},
					Date: "2021-07-19T16:42:46Z",
				},
				Message: "add unescaped file\n",
			},
		},
	}

	assert.Equal(t, dummymeta, meta)
}

func TestAPINewWikiPage(t *testing.T) {
	for _, title := range []string{
		"New page",
		"&&&&",
	} {
		defer tests.PrepareTestEnv(t)()
		username := "user2"
		session := loginUser(t, username)
		token := getTokenForLoggedInUser(t, session)

		urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/new?token=%s", username, "repo1", token)

		req := NewRequestWithJSON(t, "POST", urlStr, &api.CreateWikiPageOptions{
			Title:         title,
			ContentBase64: base64.StdEncoding.EncodeToString([]byte("Wiki page content for API unit tests")),
			Message:       "",
		})
		session.MakeRequest(t, req, http.StatusCreated)
	}
}

func TestAPIEditWikiPage(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	username := "user2"
	session := loginUser(t, username)
	token := getTokenForLoggedInUser(t, session)

	urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/Page-With-Spaced-Name?token=%s", username, "repo1", token)

	req := NewRequestWithJSON(t, "PATCH", urlStr, &api.CreateWikiPageOptions{
		Title:         "edited title",
		ContentBase64: base64.StdEncoding.EncodeToString([]byte("Edited wiki page content for API unit tests")),
		Message:       "",
	})
	session.MakeRequest(t, req, http.StatusOK)
}

func TestAPIListPageRevisions(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	username := "user2"
	session := loginUser(t, username)

	urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/revisions/Home", username, "repo1")

	req := NewRequest(t, "GET", urlStr)
	resp := session.MakeRequest(t, req, http.StatusOK)

	var revisions *api.WikiCommitList
	DecodeJSON(t, resp, &revisions)

	dummyrevisions := &api.WikiCommitList{
		WikiCommits: []*api.WikiCommit{
			{
				ID: "2c54faec6c45d31c1abfaecdab471eac6633738a",
				Author: &api.CommitUser{
					Identity: api.Identity{
						Name:  "Ethan Koenig",
						Email: "ethantkoenig@gmail.com",
					},
					Date: "2017-11-27T04:31:18Z",
				},
				Committer: &api.CommitUser{
					Identity: api.Identity{
						Name:  "Ethan Koenig",
						Email: "ethantkoenig@gmail.com",
					},
					Date: "2017-11-27T04:31:18Z",
				},
				Message: "Add Home.md\n",
			},
		},
		Count: 1,
	}

	assert.Equal(t, dummyrevisions, revisions)
}
