package core

import (
	"encoding/json"
	"fmt"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/types"
	"github.com/gocraft/web"
)

/*

  Plural Types

*/
func (*API) GetUsers(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET USERS")
}

func (*API) GetPoets(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POETS")
}

func (*API) GetPoems(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POEMS")
}

func (*API) GetIssues(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET ISSUES")
}

func (*API) GetCommittees(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET COMMITTEES")
}

/*

  User CRUD

*/
// TODO send email verification
func (a *API) RegisterUser(rw web.ResponseWriter, req *web.Request) {
	// this should be a request handler for a registration endpoint

	// get data from request

	// create registration token

	// store registration token and data in db

	// send email with a link which will have the token in the url
	// and the page you are directed to will send a POST request with
	// the provided token which will then fetch the corresponding user
	// data temporarily keyed by the token and hit the CreateUser hanlder.
}

func (a *API) CreateUser(rw web.ResponseWriter, req *web.Request) {
	var (
		user types.User
		err  error
	)

	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&user)
	if err != nil {
		a.Error("Unable to decode POST raw-data")

		// TODO send failure response to client
	}

	err = user.Validate(consts.CREATE)
	if err != nil {
		a.Error(err.Error())
	}

	// TODO once we implement email verification tokens, we will be expecting
	// a registration token in the data, which we will here use to fetch the
	// appropriate data from the db. Once we have the user, we will proceed as
	// normally below. ✲´*。.❄¨¯`*✲。❄。*。✲´*。.❄¨¯`*✲。❄。*。✲´*。.❄¨¯`*✲。❄。*。

	// insert data into db tables

	// send success response

	fmt.Println("TODO CREATE USER")
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

/*

  Poet CRD

*/
func (*API) CreatePoet(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO Create POET")
}

func (*API) GetPoet(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POET")
}

func (*API) UpdatePoet(rw web.ResponseWriter, req *web.Request) {
	// we can only update certain fields, like description and avatar
	fmt.Println("TODO UPDATE POET")
}

func (*API) DeletePoet(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO DELETE POET")
}

/*

  Poem Read

*/
func (*API) GetPoem(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POEM")
}

/*

  Issue Read

*/
func (*API) GetIssue(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET ISSUE")
}

/*

  Committee Read

*/
func (*API) GetCommittee(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET COMMITTEE")
}
