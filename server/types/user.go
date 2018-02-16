package types

import (
	"database/sql"
	"fmt"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/utils"
	_ "github.com/lib/pq"
)

type User struct {
	Id       string  `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Email    string  `json:"email"`
	Poets    []*Poet `json:"poets"`
}

func (u *User) Validate(action string) error {
	// make sure id, if not an empty string, is a uuid
	if !utils.IsValidUUIDV4(u.Id) && u.Id != "" {
		return fmt.Errorf("User Id must be a valid uuid, given %s", u.Id)
	}

	// perform validation on a per action basis
	switch action {
	case "CREATE":
	case "UPDATE":
	default:
		// only ensure that the id is present
		// this aplies to the READ and DELETE cases
	}

	return nil
}

/*
   db methods
*/

// package level globals for storing prepared sql statements
var (
	userCreateStmt  *sql.Stmt
	userReadStmt    *sql.Stmt
	userReadAllStmt *sql.Stmt
	userDeleteStmt  *sql.Stmt
)

func (*User) CreateTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS users (
		          id UUID NOT NULL UNIQUE,
                          username VARCHAR(255) NOT NULL UNIQUE,
                          password VARCHAR(255) NOT NULL,
                          email VARCHAR(255) NOT NULL UNIQUE,
		          PRIMARY KEY (id)
	)`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) CreateUser(id string, db *sql.DB) error {
	var (
		rows *sql.Rows
		err  error
	)

	// we assume that all validation/sanitization has already been called

	// assign id
	u.Id = id

	// prepare statement if not already done so.
	if userCreateStmt == nil {
		// create statement
		stmt := `INSERT INTO users (
                           id, username, password, email
                         ) VALUES ($1, $2, $3, $4)`
		userCreateStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	rows, err = userCreateStmt.Query(u.Id, u.Username, u.Password, u.Email)
	if err != nil {
		return err
	}

	// read created result
	for rows.Next() {
		rows.Scan(u.Id, u.Username, u.Password, u.Email)
	}

	return nil
}

func (u *User) ReadUser(db *sql.DB) error {
	var (
		err error
	)

	// prepare statement if not already done so.
	if userReadStmt == nil {
		// read statement
		stmt := `SELECT * FROM users WHERE id = ?`
		userReadStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	// make sure user Id is actually populated

	// run prepared query over arguments
	// NOTE: we are not joining from the poets tables
	rows, err := userReadStmt.Query(u.Id)
	if err != nil {
		return err
	}

	// decode results into user struct
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Username, &u.Password, &u.Email)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	fmt.Println(u)

	return nil
}

func (u *User) UpdateUser(db *sql.DB) error {

	return nil
}

func (u *User) DeleteUser(db *sql.DB) error {
	var (
		err error
	)

	// prepare statement if not already done so.
	if userDeleteStmt == nil {
		// delete statement
		stmt := `DELETE FROM users WHERE id = ?`
		userDeleteStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadUsers(db *sql.DB) ([]*User, error) {
	var (
		users = []*User{}
		err   error
	)

	// prepare statement if not already done so.
	if userReadAllStmt == nil {
		// readAll statement
		// TODO pagination
		stmt := `SELECT * FROM users`
		userReadAllStmt, err = db.Prepare(stmt)
		if err != nil {
			return users, nil
		}
	}

	return users, nil
}
