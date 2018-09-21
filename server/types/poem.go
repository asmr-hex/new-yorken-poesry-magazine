package types

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/utils"
	"github.com/lib/pq"
)

type Poem struct {
	Id      string
	Title   string
	Date    time.Time
	Author  *Poet
	Content string
	Issue   *Issue
	Score   float64 // score assigned by committee
	Likes   int     // number of users who liked this poem
}

// this struct is strictly for extracting possibly null valued
// fields from the database -___-
// we will only use this struct if we are OUTER JOINING poets on
// another table (e.g. users, since some users might not have poets)
// TODO (cw|9.20.2018) figure out a better way to do this...
type PoemNullable struct {
	Id      sql.NullString
	Title   sql.NullString
	Date    pq.NullTime
	Content sql.NullString
	Score   sql.NullFloat64 // score assigned by committee
}

func (pn *PoemNullable) Convert() *Poem {
	return &Poem{
		Id:      pn.Id.String,
		Title:   pn.Title.String,
		Date:    pn.Date.Time,
		Content: pn.Content.String,
		Score:   pn.Score.Float64,
	}
}

func (p *Poem) Validate(action string) error {
	// make sure id, if not empty string, is a uuid
	if !utils.IsValidUUIDV4(p.Id) && p.Id != "" {
		return fmt.Errorf("Poem Id must be a valid uuid, given %s", p.Id)
	}

	switch action {
	case consts.CREATE:
		if p.Id == "" {
			return fmt.Errorf("No id provided.")
		}
	}

	return nil
}

/*
   db methods
*/

var (
	poemCreateStmt *sql.Stmt
)

func CreatePoemsTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS poems (
		          id UUID NOT NULL UNIQUE,
                          title VARCHAR(255) NOT NULL,
                          date TIMESTAMP WITH TIME ZONE,
                          author UUID REFERENCES poets NOT NULL,
                          content TEXT NOT NULL,
                          issue UUID REFERENCES issues NOT NULL,
                          score NUMERIC(1) DEFAULT 0 CONSTRAINT normalized CHECK (score >= 0 and score <= 1),
		          PRIMARY KEY (id)
	)`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}

/*
   poem methods
*/
func (p *Poem) Create(db *sql.DB) error {
	var (
		err error
	)

	// before doing anything, lets validate this...
	err = p.Validate(consts.CREATE)
	if err != nil {
		return err
	}

	if poemCreateStmt == nil {
		stmt := `
                    INSERT INTO poems (
                      id, title, date, author, content, issue, score
                    ) VALUES ($1, $2, $3, $4, $5, $6, $7)
                `
		poemCreateStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	_, err = poemCreateStmt.Exec(
		p.Id,
		p.Title,
		p.Date,
		p.Author.Id,
		p.Content,
		p.Issue.Id,
		p.Score,
	)
	if err != nil {
		return err
	}

	return nil
}
