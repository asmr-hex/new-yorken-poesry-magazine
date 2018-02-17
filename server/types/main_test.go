package types

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/env"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

// package level globals for all test suites
var (
	testDB *sql.DB
)

// package level test entry point
func TestMain(m *testing.M) {
	var (
		retCode int
		err     error
	)

	// setup global test state
	err = setup()
	if err != nil {
		panic(err)
	}

	// run all suites
	retCode = m.Run()

	// teardown global test state
	err = teardown()
	if err != nil {
		panic(err)
	}

	os.Exit(retCode)
}

func setup() error {
	var (
		conf *env.TestConfig
		err  error
	)

	// get test config
	conf = env.NewTestConfig()

	// construct database info string required for connection
	dbInfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.DB.Host,
		conf.DB.Username,
		conf.DB.Password,
		conf.DB.Name,
		"disable",
	)

	// open a connection to the database
	testDB, err = sql.Open(env.DB_DRIVER, dbInfo)
	if err != nil {
		return err
	}

	// ping open db to verify the connection has been established.
	// otherwise (╥﹏╥)
	err = testDB.Ping()
	if err != nil {
		return err
	}

	return nil
}

func teardown() error {
	// destroy test database tables
	_, err := testDB.Exec(`DROP TABLE IF EXISTS users CASCADE;
                               DROP TABLE IF EXISTS poets CASCADE;
                               DROP TABLE IF EXISTS poems CASCADE;
                               DROP TABLE IF EXISTS issues CASCADE;
                               DROP TABLE IF EXISTS committees CASCADE;`)
	if err != nil {
		return err
	}

	return nil
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, &UserTestSuite{db: testDB})
}

func TestPoetSuite(t *testing.T) {
	suite.Run(t, &PoetTestSuite{db: testDB})
}
