package types

import (
	"database/sql"

	"github.com/stretchr/testify/suite"
)

type IssueCommitteeMembershipTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *IssueCommitteeMembershipTestSuite) SetupSuite() {
	// create users table (referenced by poets)
	err := CreateUsersTable(s.db)
	if err != nil {
		panic(err)
	}

	// create poets table
	err = CreatePoetsTable(s.db)
	if err != nil {
		panic(err)
	}

	// create issues table
	err = CreateIssuesTable(s.db)
	if err != nil {
		panic(err)
	}

}

func (s *IssueCommitteeMembershipTestSuite) TearDownSuite() {
	_, err := s.db.Exec(`DROP TABLE IF EXISTS users CASCADE`)
	if err != nil {
		panic(err)
	}

	_, err = s.db.Exec(`DROP TABLE IF EXISTS issues CASCADE`)
	if err != nil {
		panic(err)
	}
}

func (s *IssueCommitteeMembershipTestSuite) BeforeTest(suiteName, testName string) {
	var (
		err error
	)

	switch testName {
	case "TestCreateIssueCommitteeMembershipTable":
		_, err = s.db.Exec(`DROP TABLE IF EXISTS issue_committee_membership CASCADE`)
		if err != nil {
			panic(err)
		}
	case "TestAdd":
		// setup DB
		// TODO (cw|9.13.2018) setup test DB
	}
}

func (s *IssueCommitteeMembershipTestSuite) TestCreateIssueCommitteeMembershipTable() {
	err := CreateIssueCommitteeMembershipTable(testDB)
	s.NoError(err)
}

func (s *IssueCommitteeMembershipTestSuite) TestAdd_NullPoet() {
	committeeMembership := &IssueCommitteeMembership{}
	err := committeeMembership.Add(testDB)
	s.Error(err)
}
