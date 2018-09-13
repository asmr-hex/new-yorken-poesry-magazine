package types

import (
	"database/sql"
	"fmt"
	"time"
)

type Issue struct {
	Id           string
	Date         time.Time
	Committee    []*Poet
	Contributors []*Poet
	Poems        []*Poem
	Title        string
	Description  string
}

func (u *Issue) Validate() error {

	return nil
}

/*
   db methods
*/
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

func GetUpcomingIssue(db *sql.DB) (*Issue, error) {
	var (
		issue     *Issue
		committee = []*Poet{}
		err       error
	)

	// get current issue
	err = db.QueryRow(`
                SELECT id, date, title, description
                FROM issues WHERE upcoming = true
        `).Scan(
		issue.Id,
		issue.Date,
		issue.Title,
		issue.Description,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("No current issue.")
	case err != nil:
		return nil, err
	}

	// populate the committee
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

	return issue, nil
}
