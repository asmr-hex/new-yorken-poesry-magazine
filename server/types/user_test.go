package types

import (
	"database/sql"

	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	db *sql.DB
}

// run before all tests in this suite begin
func (s *UserTestSuite) SetupSuite() {
	// create users table
	user := &User{}
	err := user.CreateTable(s.db)
	if err != nil {
		panic(err)
	}
}

// run after all tests in this suite have complete
func (s *UserTestSuite) TearDownSuite() {
	_, err := s.db.Exec(`DROP TABLE IF EXISTS users CASCADE`)
	if err != nil {
		panic(err)
	}
}

// run specific setups before specific tests
func (s *UserTestSuite) BeforeTest(suiteName, testName string) {
	var (
		err error
	)

	switch testName {
	// drop users table before create table test to see if it works.
	case "TestCreateTable":
		_, err = s.db.Exec(`DROP TABLE IF EXISTS users CASCADE`)
		if err != nil {
			panic(err)
		}
	}
}

func (s *UserTestSuite) TestCreateTable() {
	user := &User{}

	err := user.CreateTable(testDB)
	s.NoError(err)
}
