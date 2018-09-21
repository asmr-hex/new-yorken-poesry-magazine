package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/utils"
	"github.com/lib/pq"
)

type Issue struct {
	Id           string
	Volume       int64
	Date         time.Time
	Committee    []*Poet
	Contributors []*Poet
	Poems        []*Poem
	Title        string
	Description  string
	Upcoming     bool
	Likes        int // number of likes this issue has
}

// this struct is strictly for extracting possibly null valued
// fields from the database -___-
// we will only use this struct if we are OUTER JOINING poets on
// another table (e.g. users, since some users might not have poets)
// TODO (cw|9.20.2018) figure out a better way to do this...
type IssueNullable struct {
	Id          sql.NullString
	Volume      sql.NullInt64
	Date        pq.NullTime
	Title       sql.NullString
	Description sql.NullString
	Upcoming    sql.NullBool
}

func (in *IssueNullable) Convert() *Issue {
	return &Issue{
		Id:          in.Id.String,
		Volume:      in.Volume.Int64,
		Date:        in.Date.Time,
		Title:       in.Title.String,
		Description: in.Description.String,
		Upcoming:    in.Upcoming.Bool,
	}
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
                          volume SERIAL
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

func ReadIssues(db *sql.DB) ([]*Issue, error) {
	var (
		issues    = []*Issue{}
		issuesMap = map[string]*Issue{}
		err       error
	)

	// TODO (cw|9.18.2018) make this more efficient...im being lazy rn?
	// read issues and judges
	rows, err := db.Query(`
                    SELECT i.id, i.date, i.title, i.description, i.upcoming,
                           j.id, j.designer, j.name, j.birthDate, j.deathDate, j.description,
                           j.language, j.programFileName, j.parameterFileName,
                           j.parameterFileIncluded, j.path
                    FROM issues i
                    INNER JOIN issue_committee_membership m
                    ON (i.id = m.issue)
                    INNER JOIN poets j
                    ON (m.poet = j.id)
                `)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		issue := &Issue{
			Committee: []*Poet{},
		}
		judge := &Poet{Designer: &User{}}
		err = rows.Scan(
			&issue.Id,
			&issue.Date,
			&issue.Title,
			&issue.Description,
			&issue.Upcoming,
			&judge.Id,
			&judge.Designer.Id,
			&judge.Name,
			&judge.BirthDate,
			&judge.DeathDate,
			&judge.Description,
			&judge.Language,
			&judge.ProgramFileName,
			&judge.ParameterFileName,
			&judge.ParameterFileIncluded,
			&judge.Path,
		)
		if err != nil {
			return nil, err
		}

		if len(issues) != 0 && issue.Id == issues[len(issues)-1].Id {
			// consolidate poets into one slice according to user
			issuesMap[issue.Id].Committee = append(
				issuesMap[issue.Id].Committee,
				judge,
			)
		} else {
			issue.Committee = []*Poet{judge}
			issuesMap[issue.Id] = issue
		}
		issues = append(issues, issue)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	issuesIds := []string{}
	// read contributors
	rows, err = db.Query(`
                    SELECT i.id,
                           c.id, c.designer, c.name, c.birthDate, c.deathDate, c.description,
                           c.language, c.programFileName, c.parameterFileName,
                           c.parameterFileIncluded, c.path
                    FROM issues i
                    INNER JOIN issue_contributions ctr
                    ON (i.id = ctr.issue)
                    INNER JOIN poets c
                    ON (ctr.poet = c.id)
                `)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var issueId string
		contributor := &Poet{Designer: &User{}}
		err = rows.Scan(
			&issueId,
			&contributor.Id,
			&contributor.Designer.Id,
			&contributor.Name,
			&contributor.BirthDate,
			&contributor.DeathDate,
			&contributor.Description,
			&contributor.Language,
			&contributor.ProgramFileName,
			&contributor.ParameterFileName,
			&contributor.ParameterFileIncluded,
			&contributor.Path,
		)
		if err != nil {
			return nil, err
		}

		if len(issuesIds) != 0 && issueId == issuesIds[len(issuesIds)-1] {
			// consolidate poets into one map according to issue
			issuesMap[issueId].Contributors = append(
				issuesMap[issueId].Contributors,
				contributor,
			)
		} else {
			issuesMap[issueId].Contributors = []*Poet{contributor}
		}
		issuesIds = append(issuesIds, issueId)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	bytes, err := json.MarshalIndent(issuesMap, "    ", "")
	fmt.Println(string(bytes))

	issues = []*Issue{}
	for _, issue := range issuesMap {
		issues = append(issues, issue)
	}

	return issues, nil
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
		poet := &Poet{Designer: &User{}}
		err = rows.Scan(
			&poet.Id,
			&poet.Designer.Id,
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

	_, err = issuePublishStmt.Exec(i.Id)
	if err != nil {
		return err
	}

	return nil
}
