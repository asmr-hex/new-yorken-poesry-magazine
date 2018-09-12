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
                          active BOOL NOT NULL,
                          issue UUID ELEMENT REFERENCES issues NOT NULL,
                          members UUID[] ELEMENT REFERENCES poets NOT NULL,
		          PRIMARY KEY (id)
	)`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}

func ReadActiveCommittee(db *sql.DB) (*ReviewCommittee, error) {
	// var (
	// 	committee *ReviewCommittee
	// 	err       error
	// )

	// TODO (cw|9.12.2018) implement me!

	return nil, nil
}
