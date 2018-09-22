package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/types"
	"github.com/gocraft/web"
	uuid "github.com/satori/go.uuid"
)

func (a *API) GetUsers(rw web.ResponseWriter, req *web.Request) {
	var (
		users []*types.User
		err   error
	)

	users, err = types.ReadUsers(a.db)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(usersJSON)
}

// TODO send email verification
func (a *API) SignUp(rw web.ResponseWriter, req *web.Request) {
	// this should be a request handler for a registration endpoint

	// get data from request

	// create registration token

	// store registration token and data in db

	// send email with a link which will have the token in the url
	// and the page you are directed to will send a POST request with
	// the provided token which will then fetch the corresponding user
	// data temporarily keyed by the token and hit the CreateUser hanlder.

	// TODO (cw|4.25.2018) everything below here will go into the VerifyAccount
	// handler.

	// For now this will *only* create a new user.
	a.CreateUser(rw, req)
}

// verifies that the account about to be created has been associated with and
// email address.
//
// TODO (cw|4.25.2018) get a verification key from the request param, lookup the
// user about to be created within an in memory cache, if it exists then proceed,
// if not, then this is an unverified email address.
//
func (a *API) VerifyAccount(rw web.ResponseWriter, req *web.Request) {
	// TODO once we implement email verification tokens, we will be expecting
	// a registration token in the data, which we will here use to fetch the
	// appropriate data from the db. Once we have the user, we will proceed as
	// normally below. ✲´*。.❄¨¯`*✲。❄。*。✲´*。.❄¨¯`*✲。❄。*。✲´*。.❄¨¯`*✲。❄。*。
}

// creates a user account.
//
// this will invariably be done after signup and email verification has occured.
//
// upon creation, a user is immediately logged in. consequently, CreateUser should
// never write to the response body so Login is able to.
//
// TODO (cw|4.25.2018) currently there is no distinction between signup and create
// user, but eventually signup should required email verification. When this occurs,
// this handler should no longer be a top-level handler, but called from the
// VerifyAccount func which will pass in the relevant info about the user about to
// be created (that info will not be in the request but cached in memory on SignUp).
// For that reason, we will eventually need to change the signature of this function
// to accept a User struct and a ResponseWriter (we no longer need the Request).
//
func (a *API) CreateUser(rw web.ResponseWriter, req *web.Request) {
	var (
		user types.User
		err  error
	)

	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&user)
	if err != nil {
		a.Error("Unable to decode POST raw-data: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	err = user.Validate(consts.CREATE)
	if err != nil {
		a.Error("User Error: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	// assign user id
	user.Id = uuid.NewV4().String()

	// insert data into db tables
	err = user.Create(a.db)
	if err != nil {
		a.Error(err.Error())

		// send response
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	a.Info("user %s successfully created!", user.Username)

	// after user is created we can then immediately log the user in
	a.login(&user, rw)
}

// handles login requests.
//
// this handler services requests from both the public api and the
// webapp dashboard. the inbound request should be have `Method: POST`
// with the credentials data json encoded. Here is an example curl,
//
// curl \
//  -X POST \
//  -H "Content-Type: application/json" \
//  -d '{"username": "percival", "password": "phantoms moving mistily"}'
//  https://poem.cool/dashboard/login
//  ( or https://poem.cool/api/v1/login)
//
// upon successful authentication, a unique session token (with an
// expiration date) is assigned to the user if not done so already. This
// token is provided in the response as a cookie, session_token.
//
// finally, the existing/authenticated user is logged in.
//
func (a *API) Login(rw web.ResponseWriter, req *web.Request) {
	var (
		user types.User
		err  error
	)

	defer req.Body.Close()

	// decode body data into a User struct
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&user)
	if err != nil {
		a.Error("Unable to decode POST raw-data: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	err = user.Validate(consts.LOGIN)
	if err != nil {
		a.Error("User Error: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	// authenticate user with password
	err = user.Authenticate(a.db)
	if err != nil {
		a.Error("User Error: %s", err.Error())

		// return response
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	a.login(&user, rw)
}

// logs in an existing/authenticated user.
//
// given a User and a ResponseWriter, fetch a session token for this user,
// write the session token and a cookie in the response, write the user json
// data to the response body.
//
// TODO (cw|4.25.2018) if the request is not from a browser
// (i.e. User-Agent is blank or something), the we should include the
// session token within the response payload?
//
func (a *API) login(user *types.User, rw web.ResponseWriter) {
	var (
		err error
	)

	// get session token
	sessionToken := a.Sessions.GetTokenByUser(user.Id)

	// set the sessionToken within a response cookie
	http.SetCookie(rw, &http.Cookie{
		Name:  SESSION_TOKEN_COOKIE_NAME,
		Value: sessionToken,
	})

	// read full user data
	err = user.Read(a.db)
	if err != nil {
		a.Error(err.Error())

		// return response
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	user.Sanitize()

	// write json encoded data into response
	err = json.NewEncoder(rw).Encode(user)
	if err != nil {
		a.Error(err.Error())

		// return response
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	a.Info("user %s successfully logged in!", user.Username)
}

func (a *API) GetUser(rw web.ResponseWriter, req *web.Request) {
	var (
		err error
	)

	// extracting the id path param
	id := req.PathParams[API_ID_PATH_PARAM]

	// assigning said id to id of user struct
	user := &types.User{Id: id}
	err = user.Validate(consts.READ)
	if err != nil {
		a.Error(err.Error())
	}

	// invoke read
	err = user.Read(a.db)
	if err != nil {
		a.Error(err.Error())
	}

	// send actual response back to users

	fmt.Println("TODO GET USER")
}

func (*API) UpdateUser(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO UPDATE USER")
}

func (*API) DeleteUser(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO DELETE USER")
}
