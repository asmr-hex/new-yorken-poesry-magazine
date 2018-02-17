package types

import (
	"database/sql"
	"path"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type PoetTestSuite struct {
	suite.Suite
	db *sql.DB
}

// run before all tests in this suite begin
func (s *PoetTestSuite) SetupSuite() {
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
}

// run after all tests in this suite have complete
func (s *PoetTestSuite) TearDownSuite() {
	_, err := s.db.Exec(`DROP TABLE IF EXISTS poets CASCADE`)
	if err != nil {
		panic(err)
	}
}

// run specific setups before specific tests
func (s *PoetTestSuite) BeforeTest(suiteName, testName string) {
	var (
		err error
	)

	switch testName {
	// drop users table before create table test to see if it works.
	case "TestCreateTable":
		_, err = s.db.Exec(`DROP TABLE IF EXISTS poets CASCADE`)
		if err != nil {
			panic(err)
		}
	}
}

func (s *PoetTestSuite) TestCreateTable() {
	err := CreatePoetsTable(testDB)
	s.NoError(err)
}

func (s *PoetTestSuite) TestCreatePoet() {
	userId := uuid.NewV4().String()
	poetId := uuid.NewV4().String()

	// create user
	user := &User{Username: "3jane", Password: "pwd", Email: "3j4n3@tessier.gov"}
	err := user.Create(userId, s.db)
	s.NoError(err)

	// create poet
	poet := &Poet{
		Designer:    userId,
		Name:        "wintermute",
		Description: "mutator of the immutable",
		ExecPath:    path.Join("/poets/", poetId),
	}

	err = poet.Create(poetId, s.db)
	s.NoError(err)
}

func (s *PoetTestSuite) TestReadPoet() {
	userId := uuid.NewV4().String()
	poetId := uuid.NewV4().String()

	// create user
	user := &User{Username: "hamilton", Password: "pwd", Email: "ijk@quaternion.idk"}
	err := user.Create(userId, s.db)
	s.NoError(err)

	// create poet
	poet := &Poet{
		Designer:    userId,
		Name:        "Chum of Chance",
		Description: "explorer of some other dimensionality",
		ExecPath:    path.Join("/poets/", poetId),
	}

	err = poet.Create(poetId, s.db)
	s.NoError(err)

	expectedPoet := poet
	expectedPoet.ExecPath = "" // this should not be public info

	// read poet
	poet = &Poet{Id: poetId}
	err = poet.Read(s.db)
	s.NoError(err)

	// since there isa problem with the postgres and golang time formats w.r.t.
	// timezones, we will just compoare the formtted times here and nillify the
	// times int he structs -__-
	expectedBirthDate := expectedPoet.BirthDate.Format(time.RFC3339)
	expectedDeathDate := expectedPoet.DeathDate.Format(time.RFC3339)
	birthDate := poet.BirthDate.Format(time.RFC3339)
	deathDate := poet.DeathDate.Format(time.RFC3339)

	s.EqualValues(expectedBirthDate, birthDate)
	s.EqualValues(expectedDeathDate, deathDate)

	expectedPoet.BirthDate = time.Time{}
	expectedPoet.DeathDate = time.Time{}
	poet.BirthDate = time.Time{}
	poet.DeathDate = time.Time{}

	s.EqualValues(expectedPoet, poet)
}
