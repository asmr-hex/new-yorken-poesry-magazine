package types

import (
	"database/sql"
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
                          volume SERIAL,
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
		// issues = map[string]*Issue{}
		issues = map[string]map[string]map[string]interface{}{}

		err error
	)

	// read issues and judges
	rows, err := db.Query(`
                    SELECT i.id, i.date, i.title, i.description, i.upcoming,

                           j.id, j.designer, j.name, j.birthDate, j.deathDate, j.description,
                           j.language, j.programFileName, j.parameterFileName,
                           j.parameterFileIncluded, j.path,

                           ju.id, ju.username, ju.email,

                           c.id, c.designer, c.name, c.birthDate, c.deathDate, c.description,
                           c.language, c.programFileName, c.parameterFileName,
                           c.parameterFileIncluded, c.path,

                           cu.id, cu.username, cu.email,

                           p.id, p.title, p.date, p.author, p.content, p.issue, p.score
                    FROM issues i

                    INNER JOIN issue_committee_membership m
                    ON (i.id = m.issue)
                    INNER JOIN poets j
                    ON (m.poet = j.id)

                    INNER JOIN users ju
                    ON (j.designer = ju.id)

                    INNER JOIN issue_contributions ctr
                    ON (i.id = ctr.issue)
                    INNER JOIN poets c
                    ON (ctr.poet = c.id)

                    INNER JOIN users cu
                    ON (c.designer = cu.id)

                    INNER JOIN poems p
                    ON (i.id = p.issue)
                `)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		issue := &Issue{
			Committee:    []*Poet{},
			Contributors: []*Poet{},
			Poems:        []*Poem{},
		}
		judgeNullable := &PoetNullable{}
		judgeDesignerNullable := &UserNullable{}
		contributorNullable := &PoetNullable{}
		contributorDesignerNullable := &UserNullable{}
		poemNullable := &PoemNullable{}
		err = rows.Scan(
			&issue.Id,
			&issue.Date,
			&issue.Title,
			&issue.Description,
			&issue.Upcoming,
			&judgeNullable.Id,
			&judgeNullable.DesignerId,
			&judgeNullable.Name,
			&judgeNullable.BirthDate,
			&judgeNullable.DeathDate,
			&judgeNullable.Description,
			&judgeNullable.Language,
			&judgeNullable.ProgramFileName,
			&judgeNullable.ParameterFileName,
			&judgeNullable.ParameterFileIncluded,
			&judgeNullable.Path,
			&judgeDesignerNullable.Id,
			&judgeDesignerNullable.Username,
			&judgeDesignerNullable.Email,
			&contributorNullable.Id,
			&contributorNullable.DesignerId,
			&contributorNullable.Name,
			&contributorNullable.BirthDate,
			&contributorNullable.DeathDate,
			&contributorNullable.Description,
			&contributorNullable.Language,
			&contributorNullable.ProgramFileName,
			&contributorNullable.ParameterFileName,
			&contributorNullable.ParameterFileIncluded,
			&contributorNullable.Path,
			&contributorDesignerNullable.Id,
			&contributorDesignerNullable.Username,
			&contributorDesignerNullable.Email,
			&poemNullable.Id,
			&poemNullable.Title,
			&poemNullable.Date,
			&poemNullable.AuthorId,
			&poemNullable.Content,
			&poemNullable.IssueId,
			&poemNullable.Score,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := issues[issue.Id]; !ok {
			issues[issue.Id] = map[string]map[string]interface{}{
				"issue":        map[string]interface{}{issue.Id: issue},
				"judges":       map[string]interface{}{},
				"contributors": map[string]interface{}{},
				"poems":        map[string]interface{}{},
			}
		}

		// oki oki, this issue has been scanned already, lets fill it in...
		// consolidate judges, contributors, and poems into slices

		// quick closure to detect if a field has already been scanned...
		// im sorry i know this is confusing.... :(((((())))))
		alreadyScanned := func(field, fieldId string) bool {
			_, ok := issues[issue.Id][field][fieldId]
			return ok
		}

		// insert judge, contributor, and poem into slices if they aren't null
		if judgeNullable.Id.Valid {
			// has this judge already been scanned?
			if ok := alreadyScanned("judges", judgeNullable.Id.String); !ok {
				designer := judgeDesignerNullable.Convert()
				judge := judgeNullable.Convert()
				judge.Designer = designer
				issues[issue.Id]["judges"][judge.Id] = judge
			}
		}

		if contributorNullable.Id.Valid {
			// has this contributor already been scanned?
			if ok := alreadyScanned("contributors", contributorNullable.Id.String); !ok {
				designer := contributorDesignerNullable.Convert()
				contributor := contributorNullable.Convert()
				contributor.Designer = designer
				issues[issue.Id]["contributors"][contributor.Id] = contributor
			}
		}

		if poemNullable.Id.Valid {
			// has this contributor already been scanned?
			if ok := alreadyScanned("poems", poemNullable.Id.String); !ok {
				poem := poemNullable.Convert()
				issues[issue.Id]["poems"][poem.Id] = poem
			}
		}

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	issuesSlice := []*Issue{}
	for id, typeMap := range issues {
		issue := typeMap["issue"][id].(*Issue)

		// append judges
		for _, judgeInt := range typeMap["judges"] {
			issue.Committee = append(
				issue.Committee,
				judgeInt.(*Poet),
			)
		}

		// append contributors
		for _, contributorInt := range typeMap["contributors"] {
			issue.Contributors = append(
				issue.Contributors,
				contributorInt.(*Poet),
			)
		}

		// append poems
		for _, poemInt := range typeMap["poems"] {
			issue.Poems = append(
				issue.Poems,
				poemInt.(*Poem),
			)
		}

		issuesSlice = append(issuesSlice, issue)
	}

	return issuesSlice, nil
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
