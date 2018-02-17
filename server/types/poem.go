package types

import (
	"database/sql"
	"time"
)

type Poem struct {
	Id      string
	Date    time.Time
	Author  *Poet
	Content string
	Issue   *Issue
	Likes   int     // number of likes from readers
	Score   float64 // score assigned by committee
}

func (p *Poem) Validate() error {
	return nil
}

/*
   db methods
*/
func (*Poem) CreateTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS poems (
		          id UUID NOT NULL UNIQUE,
                          date TIMESTAMP WITH TIME ZONE,
                          author REFERENCES poets,
                          content TEXT NOT NULL,
                          issue UUID REFERENCES issues,
                          likes INT DEFAULT 0 CONSTRAINT gtzero CHECK (likes > 0),
                          score NUMERIC(1) DEFAULT 0 CONSTRAINT normalized CHECK (score >= 0 and score <= 1),
		          PRIMARY KEY (id)
	)`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}
