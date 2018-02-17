package core

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gocraft/web"
)

const (
	API_VERSION       = "v1"
	API_DASHBOARD     = "dashboard"
	API_ID_PATH_PARAM = "id"
)

type API struct {
	*Logger
	Version string
	Router  *web.Router
	db      *sql.DB
}

func NewAPI(db *sql.DB) *API {
	api := API{
		Logger:  NewLogger(os.Stdout),
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
		Get(pubPrefix+"/users", a.GetUsers).
		Get(pubPrefix+"/poets", a.GetPoets).
		Get(pubPrefix+"/poems", a.GetPoems).
		Get(pubPrefix+"/issues", a.GetIssues).
		Get(pubPrefix+"/committees", a.GetCommittees).

		// User CRUD
		Post(pubPrefix+"/user", a.CreateUser).
		Get(pubPrefix+"/user/:"+API_ID_PATH_PARAM, a.GetUser).
		Put(pubPrefix+"/user/:"+API_ID_PATH_PARAM, a.UpdateUser).
		Delete(pubPrefix+"/user/:"+API_ID_PATH_PARAM, a.DeleteUser).

		// Poet CRD
		Post(pubPrefix+"/poet", a.CreatePoet).
		Get(pubPrefix+"/poet/:"+API_ID_PATH_PARAM, a.GetPoet).
		Put(pubPrefix+"/poet/:"+API_ID_PATH_PARAM, a.UpdatePoet).
		Delete(pubPrefix+"/poet/:"+API_ID_PATH_PARAM, a.DeletePoet).

		// Poem R (Poems can only be read via the API)
		Get(pubPrefix+"/poem/:"+API_ID_PATH_PARAM, a.GetPoem).

		// Issue R (Issues can only be read via the API)
		Get(pubPrefix+"/issue/:"+API_ID_PATH_PARAM, a.GetIssue).

		// Committee R (Committees can only be read via the API)
		Get(pubPrefix+"/committee/:"+API_ID_PATH_PARAM, a.GetCommittee).

		// === Dashboard API ===

		// TODO
		Get(dashPrefix+"/user/:"+API_ID_PATH_PARAM, a.GetUser)
}
