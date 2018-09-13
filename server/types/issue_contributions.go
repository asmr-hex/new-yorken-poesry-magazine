package types

import (
	"database/sql"
	"fmt"
)

type IssueContributions struct {
	Id    string
	Poet  *Poet
	Issue *Issue
}

func CreateIssueContributionsTable(db *sql.DB) error {
	mkTableStmt := ` CREATE TABLE IF NOT EXISTS issue_contributions (
                           poet UUID REFERENCES poets NOT NULL,
                           issue UUID REFERENCES issues NOT NULL,
                           PRIMARY KEY (poet, issue)
        )`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}

var (
	addIssueContributionsStmt *sql.Stmt
)

func (i *IssueContributions) Add(db *sql.DB) error {
	var (
		err error
	)

	// ensure that each are not nil
	if i.Poet == nil || i.Issue == nil {
		return fmt.Errorf("unable to add issue contribution (null poet or issue)")
	}

	if i.Poet.Id == "" || i.Issue.Id == "" {
		return fmt.Errorf("unable to add issue contribution (invalid poet or issue id)")
	}

	if addIssueContributionsStmt == nil {
		stmt := `
                    INSERT INTO issue_contributions (
                      poet, issue
                    ) VALUES ($1, $2)
                `
		addIssueContributionsStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	_, err = addIssueContributionsStmt.Exec(i.Poet.Id, i.Issue.Id)
	if err != nil {
		return err
	}

	return nil
}
