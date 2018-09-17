package types

import (
	"database/sql"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	db *sql.DB
}

// run before all tests in this suite begin
func (s *UserTestSuite) SetupSuite() {
	// create users table
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
	case "TestReadAllUsers":
		_, err = s.db.Exec(`DROP TABLE IF EXISTS users CASCADE`)
		if err != nil {
			panic(err)
		}
		err := CreateUsersTable(testDB)
		if err != nil {
			panic(err)
		}
	}
}

func (s *UserTestSuite) TestCreateTable() {
	err := CreateUsersTable(testDB)
	s.NoError(err)
}

func (s *UserTestSuite) TestCreateUser() {
	id := uuid.NewV4().String()
	user := &User{Id: id, Username: "tlon", Password: "axaxaxax", Email: "hr@worst.nightmare"}
	expectedResult := &User{
		Id:       id,
		Username: "tlon",
		Password: "axaxaxax",
		Email:    "hr@worst.nightmare",
	}

	err := user.Create(testDB)
	s.NoError(err)
	s.EqualValues(user, expectedResult)
}

func (s *UserTestSuite) TestReadUser() {
	id := uuid.NewV4().String()
	expectedUser := &User{
		Id:       id,
		Username: "dagon",
		Password: "bl4ckr33f",
		Email:    "gasp@unknowable.horror",
		Poets:    []*Poet{},
	}

	// create expected user in db
	err := expectedUser.Create(testDB)
	s.NoError(err)

	user := &User{Id: id}
	err = user.Read(testDB)
	s.NoError(err)

	// ignore passwords for now: TODO test salt
	expectedUser.Password = ""
	user.Password = ""

	s.EqualValues(expectedUser, user)
}

func (s *UserTestSuite) TestReadUser_NonExistent() {
	user := &User{Id: uuid.NewV4().String()}
	err := user.Read(testDB)
	s.Error(err)
}

func (s *UserTestSuite) TestDeleteUser() {
	id := uuid.NewV4().String()
	// create a user
	user := &User{
		Id:       id,
		Username: "colonel_buendias",
		Password: "g0ld3nf15h",
		Email:    "rogue@macondo.gov",
	}
	err := user.Create(testDB)
	s.NoError(err)

	// delete a user
	err = user.Delete(testDB)
	s.NoError(err)

	// read a user (shouldn't exist)
	err = user.Read(testDB)
	s.Error(err)
}

func (s *UserTestSuite) TestReadAllUsers() {
	var err error

	// create three users
	ids := []string{uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String()}
	expectedUsers := []*User{
		&User{Id: ids[0], Username: "a", Password: "passwd", Email: "a@idk.org"},
		&User{Id: ids[1], Username: "b", Password: "passwd", Email: "b@idk.org"},
		&User{Id: ids[2], Username: "c", Password: "passwd", Email: "c@idk.org"},
	}

	for i := 0; i < len(ids); i++ {
		err = expectedUsers[i].Create(testDB)
		s.NoError(err)

		// since we do not read passwords of users, we set them to empty string
		expectedUsers[i].Password = ""
	}

	// read all users
	users, err := ReadUsers(testDB)
	s.NoError(err)

	s.EqualValues(expectedUsers, users)
}
