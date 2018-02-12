package core

import (
	"fmt"

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
func (*API) CreateUser(rw web.ResponseWriter, req *web.Request) {
	// TODO send email verification

	// extract username, password, email, etc from form-data

	// validate all provided form fields

	// insert data into db tables

	// send success response

	fmt.Println("TODO CREATE USER")
}

func (*API) GetUser(rw web.ResponseWriter, req *web.Request) {
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
