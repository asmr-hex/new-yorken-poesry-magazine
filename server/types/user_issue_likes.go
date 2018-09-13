package types

import (
	"database/sql"
	"fmt"
)

type UserIssueLikes struct {
	Id    string
	User  *User
	Issue *Issue
}

func CreateUserIssueLikesTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS user_issue_likes (
                          usr UUID REFERENCES users NOT NULL,
                          issue UUID REFERENCES issues NOT NULL,
                          PRIMARY KEY (usr, issue)
        )`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}

var (
	addUserIssueLikeStmt *sql.Stmt
)

func (i *UserIssueLikes) Add(db *sql.DB) error {
	var (
		err error
	)

	// ensure that each are not nil
	if i.User == nil || i.Issue == nil {
		return fmt.Errorf("unable to add user issue like (null user or issue)")
	}

	if i.User.Id == "" || i.Issue.Id == "" {
		return fmt.Errorf("unable to add user issue like (invalid user or issue id)")
	}

	if addUserIssueLikeStmt == nil {
		stmt := `
                    INSERT INTO user_issue_likes (
                      usr, issue
                    ) VALUES ($1, $2)
                `
		addUserIssueLikeStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	_, err = addUserIssueLikeStmt.Exec(i.User.Id, i.Issue.Id)
	if err != nil {
		return err
	}

	return nil
}
