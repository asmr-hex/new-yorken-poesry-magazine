package core

import (
	"fmt"

	"github.com/gocraft/web"
)

const (
	VERSION   = "v1"
	DASHBOARD = "dashboard"
)

type API struct {
	Version string
	Router  *web.Router
}

func NewAPI() *API {
	api := API{Version: VERSION}

	api.BuildRouter()

	return &api
}

func (a *API) BuildRouter() {
	var (
		prefix     = "/api"
		pubPrefix  = fmt.Sprintf("%s/%s", prefix, a.Version)
		dashPrefix = fmt.Sprintf("%s/%s", prefix, DASHBOARD)
	)

	a.Router = web.New(*a).
		// TODO include middlewares for authorization etc.

		// === Public API ===

		// Plural type Reads
		Get(pubPrefix+"/users", (*API).GetUsers).
		Get(pubPrefix+"/poets", (*API).GetPoets).
		Get(pubPrefix+"/poems", (*API).GetPoems).
		Get(pubPrefix+"/issues", (*API).GetIssues).
		Get(pubPrefix+"/committees", (*API).GetCommittees).

		// User CRUD
		Post(pubPrefix+"/user", (*API).CreateUser).
		Get(pubPrefix+"/user/id", (*API).GetUser).
		Put(pubPrefix+"/user/id", (*API).UpdateUser).
		Delete(pubPrefix+"/user/id", (*API).DeleteUser).

		// Poet CRD
		Post(pubPrefix+"/poet", (*API).CreatePoet).
		Get(pubPrefix+"/poet/id", (*API).GetPoet).
		Put(pubPrefix+"/poet/id", (*API).UpdatePoet).
		Delete(pubPrefix+"/poet/id", (*API).DeletePoet).

		// Poem R (Poems can only be read via the API)
		Get(pubPrefix+"/poem/id", (*API).GetPoem).

		// Issue R (Issues can only be read via the API)
		Get(pubPrefix+"/issue/id", (*API).GetIssue).

		// Committee R (Committees can only be read via the API)
		Get(pubPrefix+"/committee/id", (*API).GetCommittee).

		// === Dashboard API ===

		// TODO
		Get(dashPrefix+"/user/id", (*API).GetUser)
}
