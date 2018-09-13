package types

import (
	"database/sql"
	"fmt"
)

type UserPoetLikes struct {
	Id   string
	User *User
	Poet *Poet
}

func CreateUserPoetLikesTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS user_poet_likes (
                           usr UUID REFERENCES users NOT NULL,
                           poet UUID REFERENCES poets NOT NULL,
                           PRIMARY KEY (usr, poet)
        )`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}

var (
	addUserPoetLikesStmt *sql.Stmt
)

func (i *UserPoetLikes) Add(db *sql.DB) error {
	var (
		err error
	)

	// ensure that each are not nil
	if i.User == nil || i.Poet == nil {
		return fmt.Errorf("unable to add user poet like (null user or poet)")
	}

	if i.User.Id == "" || i.Poet.Id == "" {
		return fmt.Errorf("unable to add user poet like (invalid user or poet id)")
	}

	if addUserPoetLikesStmt == nil {
		stmt := `
                    INSERT INTO user_poet_likes (
                      usr, poet
                    ) VALUES ($1, $2)
                `
		addUserPoetLikesStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	_, err = addUserPoetLikesStmt.Exec(i.User.Id, i.Poet.Id)
	if err != nil {
		return err
	}

	return nil
}
