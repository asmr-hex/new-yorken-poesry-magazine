package core

import (
	"database/sql"
	"fmt"

	"github.com/gocraft/web"
)

const (
	API_VERSION       = "v1"
	API_DASHBOARD     = "dashboard"
	API_ID_PATH_PARAM = "id"
)

type API struct {
	Version string
	Router  *web.Router
	db      *sql.DB
}

func NewAPI(db *sql.DB) *API {
	api := API{
		Version: API_VERSION,
		db:      db,
	}

	api.BuildRouter()

	return &api
}

func (a *API) BuildRouter() {
	// split on two api routes, Pulic (versioned) and Dashboard
	var (
		prefix     = "/api"
		pubPrefix  = fmt.Sprintf("%s/%s", prefix, a.Version)
		dashPrefix = fmt.Sprintf("%s/%s", prefix, API_DASHBOARD)
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
		Get(pubPrefix+"/user/:"+API_ID_PATH_PARAM, (*API).GetUser).
		Put(pubPrefix+"/user/:"+API_ID_PATH_PARAM, (*API).UpdateUser).
		Delete(pubPrefix+"/user/:"+API_ID_PATH_PARAM, (*API).DeleteUser).

		// Poet CRD
		Post(pubPrefix+"/poet", (*API).CreatePoet).
		Get(pubPrefix+"/poet/:"+API_ID_PATH_PARAM, (*API).GetPoet).
		Put(pubPrefix+"/poet/:"+API_ID_PATH_PARAM, (*API).UpdatePoet).
		Delete(pubPrefix+"/poet/:"+API_ID_PATH_PARAM, (*API).DeletePoet).

		// Poem R (Poems can only be read via the API)
		Get(pubPrefix+"/poem/:"+API_ID_PATH_PARAM, (*API).GetPoem).

		// Issue R (Issues can only be read via the API)
		Get(pubPrefix+"/issue/:"+API_ID_PATH_PARAM, (*API).GetIssue).

		// Committee R (Committees can only be read via the API)
		Get(pubPrefix+"/committee/:"+API_ID_PATH_PARAM, (*API).GetCommittee).

		// === Dashboard API ===

		// TODO
		Get(dashPrefix+"/user/:"+API_ID_PATH_PARAM, (*API).GetUser)
}
