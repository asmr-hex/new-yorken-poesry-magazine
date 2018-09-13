package types

import (
	"database/sql"
	"fmt"
)

type IssueCommitteeMembership struct {
	Id    string
	Poet  *Poet
	Issue *Issue
}

func CreateIssueCommitteeMembershipTable(db *sql.DB) error {
	mkTableStmt := ` CREATE TABLE IF NOT EXISTS issue_committee_membership (
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
	addCommitteeMemberStmt *sql.Stmt
)

func (i *IssueCommitteeMembership) Add(db *sql.DB) error {
	var (
		err error
	)

	// ensure that each are not nil
	if i.Poet == nil || i.Issue == nil {
		return fmt.Errorf("unable to add committee member (null poet or issue)")
	}

	if i.Poet.Id == "" || i.Issue.Id == "" {
		return fmt.Errorf("unable to add committee member (invalid poet or issue id)")
	}

	if addCommitteeMemberStmt == nil {
		stmt := `
                    INSERT INTO issue_committee_membership (
                      poet, issue
                    ) VALUES ($1, $2)
                `
		addCommitteeMemberStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	_, err = addCommitteeMemberStmt.Exec(i.Poet.Id, i.Issue.Id)
	if err != nil {
		return err
	}

	return nil
}
