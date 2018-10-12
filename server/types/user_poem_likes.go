package types

import (
	"database/sql"
	"fmt"
)

type UserPoemLikes struct {
	Id   string
	User *User
	Poem *Poem
}

func CreateUserPoemLikesTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS user_poem_likes (
                           usr UUID REFERENCES users NOT NULL,
                           poem UUID REFERENCES poems NOT NULL,
                           PRIMARY KEY (usr, poem)
        )`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}

var (
	addUserPoemLikesStmt *sql.Stmt
)

func (i *UserPoemLikes) Add(db *sql.DB) error {
	var (
		err error
	)

	// ensure that each are not nil
	if i.User == nil || i.Poem == nil {
		return fmt.Errorf("unable to add user poem like (null user or poem)")
	}

	if i.User.Id == "" || i.Poem.Id == "" {
		return fmt.Errorf("unable to add user poem like (invalid user or poem id)")
	}

	if addUserPoemLikesStmt == nil {
		stmt := `
                    INSERT INTO user_poem_likes (
                      usr, poem
                    ) VALUES ($1, $2)
                `
		addUserPoemLikesStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	_, err = addUserPoemLikesStmt.Exec(i.User.Id, i.Poem.Id)
	if err != nil {
		return err
	}

	return nil
}
