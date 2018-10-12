package core

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/env"
	"github.com/gocraft/web"
)

const (
	API_VERSION                 = "v1"
	API_DASHBOARD               = "dashboard"
	API_ID_PATH_PARAM           = "id"
	API_ISSUE_VOLUME_PATH_PARAM = "volume"
	API_VERIFY_PATH_PARAM       = "verify"
)

type API struct {
	*Logger
	Config   *env.Config
	Version  string
	Router   *web.Router
	Sessions *Sessions
	db       *sql.DB
}

func NewAPI(config *env.Config, db *sql.DB) *API {
	api := API{
		Logger:  NewLogger(os.Stdout),
		Config:  config,
		Version: API_VERSION,
		db:      db,
	}

	api.BuildRouter()

	api.Sessions = NewSessions(time.Minute * 30) // TODO put in config

	return &api
}

func (a *API) BuildRouter() {
	// split on two api routes, Pulic (versioned) and Dashboard
	var (
		pubPrefix  = fmt.Sprintf("api/%s", a.Version)
		dashPrefix = fmt.Sprintf("dashboard")
	)

	a.Router = web.New(*a).
		// === Middleware ===

		Middleware((*API).Validate).
		Middleware((*API).Authorize).

		// === Public API ===

		// User Register /Login
		Post(pubPrefix+"/register", a.CreateUser).
		Post(pubPrefix+"/login", a.Login).

		// Plural type Reads
		Get(pubPrefix+"/users", a.GetUsers).
		Get(pubPrefix+"/poets", a.GetPoets).
		Get(pubPrefix+"/poems", a.GetPoems).
		Get(pubPrefix+"/issues", a.GetIssues).

		// User RUD
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

		// Code by Poet Id Read
		Get(pubPrefix+"/code/:"+API_ID_PATH_PARAM, a.GetPoetCode).

		// Supported Languages
		Get(pubPrefix+"/supported-languages", a.GetSupportedLanguages).

		// Generate Poem
		Get(pubPrefix+"/poet/:"+API_ID_PATH_PARAM+"/write-poem", a.GeneratePoem).

		// Issue R (Issues can only be read via the API)
		Get(pubPrefix+"/issue/:"+API_ISSUE_VOLUME_PATH_PARAM, a.GetIssueByVolume).

		// === Dashboard API ===
		Post(dashPrefix+"/login", a.Login).
		Post(dashPrefix+"/signup", a.SignUp).
		Post(dashPrefix+"/verify/:"+API_VERIFY_PATH_PARAM, a.VerifyAccount).

		// poet endpoints
		Post(dashPrefix+"/poet", a.CreatePoet).
		Delete(dashPrefix+"/poet/:"+API_ID_PATH_PARAM, a.DeletePoet).

		// TODO
		Get(dashPrefix+"/user/:"+API_ID_PATH_PARAM, a.GetUser)

	// set serve static middleware only if in production since in dev env mode
	// we are using a hot-reloading node server for our frontend.
	// TODO (cw|9.15.2018) get rid of this, we are using nginx to serve static pages
	if a.Config.DevEnv == false {
		currentRoot, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// add static server middleware
		a.Router.Middleware(
			web.StaticMiddleware(
				path.Join(currentRoot, "client/build"),
				web.StaticOption{IndexFile: "index.html"},
			),
		)
	}
}
