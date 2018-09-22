package core

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/types"
	"github.com/gocraft/web"
)

func (a *API) GetIssues(rw web.ResponseWriter, req *web.Request) {
	var (
		issues []*types.Issue
		err    error
	)

	issues, err = types.ReadIssues(a.db)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	issuesJSON, err := json.Marshal(issues)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(issuesJSON)
}

func (a *API) GetIssueByVolume(rw web.ResponseWriter, req *web.Request) {
	var (
		issue *types.Issue = &types.Issue{}
	)

	// extracting the volume path param
	volume := req.PathParams[API_ISSUE_VOLUME_PATH_PARAM]

	if volume == "latest" {
		latestExists, err := issue.ReadLatest(a.db)
		if err != nil {
			a.Error(err.Error())

			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}

		if !latestExists {
			// crud the latest issue doesn't exist
			// but its okay-- this just means that the magazine
			// hasn't released the first issue yet.

			// TODO (cw|9.22.2018) return an error?

			return
		}
	} else {
		// parse int of volume number
		volumeInt, err := strconv.ParseInt(volume, 0, 64)
		if err != nil {
			a.Error(err.Error())

			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}

		issue.Volume = volumeInt

		err = issue.ReadByVolume(a.db)
		if err != nil {
			a.Error(err.Error())

			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	issueJSON, err := json.Marshal(issue)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(issueJSON)
}
