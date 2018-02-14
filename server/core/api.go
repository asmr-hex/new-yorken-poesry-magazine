package core

import (
	"fmt"

	"github.com/gocraft/web"
)

const (
	VERSION       = "v1"
	DASHBOARD     = "dashboard"
	ID_PATH_PARAM = "id"
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
		// === Middleware ===
		Middleware((*API).Validate).
		Middleware((*API).Authorize).

		// === Public API ===

		// Plural type Reads
		Get(pubPrefix+"/users", (*API).GetUsers).
		Get(pubPrefix+"/poets", (*API).GetPoets).
		Get(pubPrefix+"/poems", (*API).GetPoems).
		Get(pubPrefix+"/issues", (*API).GetIssues).
		Get(pubPrefix+"/committees", (*API).GetCommittees).

		// User CRUD
		Post(pubPrefix+"/user", (*API).CreateUser).
		Get(pubPrefix+"/user/:"+ID_PATH_PARAM, (*API).GetUser).
		Put(pubPrefix+"/user/:"+ID_PATH_PARAM, (*API).UpdateUser).
		Delete(pubPrefix+"/user/:"+ID_PATH_PARAM, (*API).DeleteUser).

		// Poet CRD
		Post(pubPrefix+"/poet", (*API).CreatePoet).
		Get(pubPrefix+"/poet/:"+ID_PATH_PARAM, (*API).GetPoet).
		Put(pubPrefix+"/poet/:"+ID_PATH_PARAM, (*API).UpdatePoet).
		Delete(pubPrefix+"/poet/:"+ID_PATH_PARAM, (*API).DeletePoet).

		// Poem R (Poems can only be read via the API)
		Get(pubPrefix+"/poem/:"+ID_PATH_PARAM, (*API).GetPoem).

		// Issue R (Issues can only be read via the API)
		Get(pubPrefix+"/issue/:"+ID_PATH_PARAM, (*API).GetIssue).

		// Committee R (Committees can only be read via the API)
		Get(pubPrefix+"/committee/:"+ID_PATH_PARAM, (*API).GetCommittee).

		// === Dashboard API ===

		// TODO
		Get(dashPrefix+"/user/:"+ID_PATH_PARAM, (*API).GetUser)
}
