package types

import (
	"database/sql"
	"fmt"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
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
	case consts.CREATE:
	case consts.UPDATE:
	case consts.DELETE:
		// TODO ensure that only a user can delete themselves
		fallthrough
	default:
		// only ensure that the id is present
		// this aplies to the READ and DELETE cases
		// we minimally need the Id to exist in these cases
		if u.Id == "" {
			return fmt.Errorf(
				"User Id is a required field to fulfill a %s",
				action,
			)
		}
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

func (u *User) Create(id string, db *sql.DB) error {
	var (
		err error
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

	_, err = userCreateStmt.Exec(u.Id, u.Username, u.Password, u.Email)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Read(db *sql.DB) error {
	var (
		err error
	)

	// prepare statement if not already done so.
	if userReadStmt == nil {
		// read statement
		stmt := `SELECT * FROM users WHERE id = $1`
		userReadStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	// make sure user Id is actually populated

	// run prepared query over arguments
	// NOTE: we are not joining from the poets tables
	err = userReadStmt.
		QueryRow(u.Id).
		Scan(&u.Id, &u.Username, &u.Password, &u.Email)
	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("No user with id %s", u.Id)
	case err != nil:
		return err
	}

	return nil
}

func (u *User) Update(db *sql.DB) error {
	// TODO

	return nil
}

func (u *User) Delete(db *sql.DB) error {
	var (
		err error
	)

	// prepare statement if not already done so.
	if userDeleteStmt == nil {
		// delete statement
		stmt := `DELETE FROM users WHERE id = $1`
		userDeleteStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	_, err = userDeleteStmt.Exec(u.Id)
	if err != nil {
		return err
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
		stmt := `SELECT id, username, email FROM users`
		userReadAllStmt, err = db.Prepare(stmt)
		if err != nil {
			return users, nil
		}
	}

	rows, err := userReadAllStmt.Query()
	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Email)
		if err != nil {
			return users, err
		}

		// append scanned user into list of all users
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}
