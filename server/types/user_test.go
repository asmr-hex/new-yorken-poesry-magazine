package types

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/env"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	testDB *sql.DB
)

const (
	TEST_DB_NAME = "nypm_test"
)

func TestMain(m *testing.M) {
	// setup db before calling tests
	err := setup()
	if err != nil {
		panic(err)
	}

	returnCode := m.Run()

	// tear down everything
	err = teardown()
	if err != nil {
		panic(err)
	}

	os.Exit(returnCode)
}

func setup() error {
	var (
		err error
	)

	// get config
	conf := env.NewConfig()

	// construct database info string required for connection
	dbInfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.DB.Host,
		conf.DB.Username,
		conf.DB.Password,
		TEST_DB_NAME, // eventually get config from .env file
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

func TestCreateTable(t *testing.T) {
	user := &User{}

	err := user.CreateTable(testDB)
	assert.NoError(t, err)
}

func TestNothing(t *testing.T) {

}
