package types

import (
	"database/sql"
	"time"
)

type Issue struct {
	Id           string
	Date         time.Time
	Committee    *ReviewCommittee
	Contributors []*Poet
	Poems        []*Poem
	Description  string
}

func (u *Issue) Validate() error {

	return nil
}

/*
   db methods
*/
func (*Issue) CreateTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS issues (
		          id UUID NOT NULL UNIQUE,
                          date TIMESTAMP WITH TIME ZONE NOT NULL,
                          committee UUID[] ELEMENT REFERENCES committees NOT NULL,
                          contributors UUID[] ELEMENT REFERENCES poets NOT NULL,
                          poems UUID[] ELEMENT REFERENCES poems NOT NULL,
                          description TEXT,
		          PRIMARY KEY (id)
	)`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}
