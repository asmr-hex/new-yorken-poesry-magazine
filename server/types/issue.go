package types

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/utils"
)

type Issue struct {
	Id           string
	Date         time.Time
	Committee    []*Poet
	Contributors []*Poet
	Poems        []*Poem
	Title        string
	Description  string
	Upcoming     bool
}

func (i *Issue) Validate(action string) error {
	// make sure id, if not empty string, is a uuid
	if !utils.IsValidUUIDV4(i.Id) && i.Id != "" {
		return fmt.Errorf("Issue Id must be a valid uuid, given %s", i.Id)
	}

	switch action {
	case consts.CREATE:
		if i.Id == "" {
			return fmt.Errorf("No id provided.")
		}
	}

	return nil
}

/*
   db methods
*/
var (
	issueCreateStmt  *sql.Stmt
	issuePublishStmt *sql.Stmt
)

func CreateIssuesTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS issues (
		          id UUID NOT NULL UNIQUE,
                          date TIMESTAMP WITH TIME ZONE NOT NULL,
                          title VARCHAR(255) NOT NULL,
                          description TEXT NOT NULL,
                          upcoming BOOL NOT NULL,
		          PRIMARY KEY (id)
	)`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}

func (i *Issue) Create(db *sql.DB) error {
	var (
		err error
	)

	// before doing anything, lets validate this...
	err = i.Validate(consts.CREATE)
	if err != nil {
		return err
	}

	if issueCreateStmt == nil {
		stmt := `
                    INSERT INTO issues (
                      id, date, title, description, upcoming
                    ) VALUES ($1, $2, $3, $4, $5)
                `
		issueCreateStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	_, err = issueCreateStmt.Exec(
		i.Id,
		i.Date,
		i.Title,
		i.Description,
		i.Upcoming,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetUpcomingIssue(db *sql.DB) (*Issue, error) {
	var (
		issue     Issue
		committee = []*Poet{}
		err       error
	)

	// TODO (cw|9.14.2018) make these queries better -___-

	// get current issue
	err = db.QueryRow(`
                SELECT id, date, title, description, upcoming
                FROM issues WHERE upcoming = true
        `).Scan(
		&issue.Id,
		&issue.Date,
		&issue.Title,
		&issue.Description,
		&issue.Upcoming,
	)
	switch {
	case err == sql.ErrNoRows:
		// tis means that there is no current issue, which must
		// mean that we need to make the FIRST ISSUE WAHOOOO
		return nil, nil
	case err != nil:
		return nil, err
	}

	// populate the committee (and no other joins since we haven't chosen poems yet)
	rows, err := db.Query(`
                SELECT p.id, p.designer, p.name, p.birthDate, p.deathDate, p.description, p.language, p.programFileName, p.parameterFileName, p.parameterFileIncluded, p.path
                FROM issue_committee_membership c
                INNER JOIN poets p
                ON (c.issue = $1 AND c.poet = p.id)
        `, issue.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		poet := &Poet{}
		err = rows.Scan(
			&poet.Id,
			&poet.Designer,
			&poet.Name,
			&poet.BirthDate,
			&poet.DeathDate,
			&poet.Description,
			&poet.Language,
			&poet.ProgramFileName,
			&poet.ParameterFileName,
			&poet.ParameterFileIncluded,
			&poet.Path,
		)
		if err != nil {
			return nil, err
		}

		committee = append(committee, poet)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	issue.Committee = committee

	return &issue, nil
}

func (i *Issue) Publish(db *sql.DB) error {
	var (
		err error
	)

	if issuePublishStmt == nil {
		stmt := `
                    UPDATE issues
                    SET upcoming = false
                    WHERE id = $1
                `
		issuePublishStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	_, err = poemCreateStmt.Exec(i.Id)
	if err != nil {
		return err
	}

	return nil
}
