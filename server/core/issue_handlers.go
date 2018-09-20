package core

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	a.Info("successfully read all issues")
}

func (*API) GetIssue(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET ISSUE")
}
