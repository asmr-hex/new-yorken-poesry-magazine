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
	Id       sql.NullString
	Title    sql.NullString
	Date     pq.NullTime
	AuthorId sql.NullString
	Content  sql.NullString
	IssueId  sql.NullString
	Score    sql.NullFloat64 // score assigned by committee
}

func (pn *PoemNullable) Convert() *Poem {
	return &Poem{
		Id:      pn.Id.String,
		Title:   pn.Title.String,
		Date:    pn.Date.Time,
		Author:  &Poet{Id: pn.AuthorId.String},
		Content: pn.Content.String,
		Issue:   &Issue{Id: pn.IssueId.String},
		Score:   pn.Score.Float64,
	}
}

var (
	poemReadAllStmt *sql.Stmt
)

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

func (p *Poem) Read(db *sql.DB) error {

	return nil
}

func ReadPoems(db *sql.DB) ([]*Poem, error) {
	var (
		poems []*Poem = []*Poem{}
		err   error
	)

	if poemReadAllStmt == nil {
		// readAll statement
		// TODO pagination
		stmt := `
                         SELECT p.id, title, date, content, score,

                                a.id, a.name, a.birthdate, a.deathdate,
                                a.description, a.language, a.programFileName,
                                a.parameterFileName, a.parameterFileIncluded,
                                a.path,

                                u.id, u.username, u.email
                         FROM poems p
                         INNER JOIN poets a ON (p.author = a.id)
                         INNER JOIN users u ON (a.designer = u.id)
                `
		poemReadAllStmt, err = db.Prepare(stmt)
		if err != nil {
			return nil, err
		}
	}

	rows, err := poemReadAllStmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		poem := &Poem{Author: &Poet{Designer: &User{}}}
		err = rows.Scan(
			&poem.Id,
			&poem.Title,
			&poem.Date,
			&poem.Content,
			&poem.Score,
			&poem.Author.Id,
			&poem.Author.Name,
			&poem.Author.BirthDate,
			&poem.Author.DeathDate,
			&poem.Author.Description,
			&poem.Author.Language,
			&poem.Author.ProgramFileName,
			&poem.Author.ParameterFileName,
			&poem.Author.ParameterFileIncluded,
			&poem.Author.Path,
			&poem.Author.Designer.Id,
			&poem.Author.Designer.Username,
			&poem.Author.Designer.Email,
		)
		if err != nil {
			return nil, err
		}

		// append scanned user into list of all poets
		poems = append(poems, poem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return poems, nil
}
