package types

import "database/sql"

type ReviewCommittee struct {
	Id      string
	Members []*Poet
}

func (c *ReviewCommittee) Validate() error {
	return nil
}

/*
   db methods
*/
func (*ReviewCommittee) CreateTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS committees (
		          id UUID NOT NULL UNIQUE,
                          members UUID[] ELEMENT REFERENCES poets NOT NULL,
		          PRIMARY KEY (id)
	)`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}
