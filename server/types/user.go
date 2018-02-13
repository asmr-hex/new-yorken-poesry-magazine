package types

import "database/sql"

type User struct {
	Id       string  `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Email    string  `json:"email"`
	Poets    []*Poet `json:"poets"`
}

func (u *User) Validate() error {
	// TODO

	return nil
}

/*
   db methods
*/
func (*User) CreateTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS users (
		          id UUID NOT NULL UNIQUE,
                          username VARCHAR(255) NOT NULL UNIQUE,
                          password VARCHAR(255) NOT NULL,
                          email VARCHAR(255) NOT NULL UNIQUE,
                          poets UUID[] ELEMENT REFERENCES poets,
		          PRIMARY KEY (id)
	)`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) CreateUser(db *sql.DB) error {

	return nil
}

func (u *User) ReadUser(db *sql.DB) error {

	return nil
}

func (u *User) UpdateUser(db *sql.DB) error {

	return nil
}

func (u *User) DeleteUser(db *sql.DB) error {

	return nil
}

func ReadUsers(db *sql.DB) ([]*User, error) {
	var (
		users = []*User{}
	)

	return users, nil
}

/*
   sql statement preparation
*/

// statements are used to optimize common read/write patterns to the db
type userStmts struct {
	prepared bool
	create   *sql.Stmt
	read     *sql.Stmt
	readAll  *sql.Stmt
	update   *sql.Stmt
	delete   *sql.Stmt
}

// create a package scoped global to store these statements ヾ(_ _。）
var (
	usrStmts = &userStmts{prepared: false}
)

// prepares statments. This function is idempotent s.t. it will only prepare statements
// if they have not yet been prepared and do nothing otherwise.
func (s *userStmts) prepareStmts(db *sql.DB) error {
	if s.prepared {
		// okidoki, stmts have already been prepared ^-^
		return nil
	}

	//

	return nil
}
