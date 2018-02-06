package core

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	DB_DRIVER = "postgres"
)

type Config struct {
	User     string
	Password string
	DBName   string
}

type Platform struct {
	config *Config
	db     *sql.DB
}

func (p *Platform) Connect() {
	var (
		err    error
		dbInfo string
	)

	dbInfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		p.config.User, p.config.Password, p.config.DBName, "disable")

	p.db, err = sql.Open(DB_DRIVER, dbInfo)
	if err != nil {
		panic(err)
	}

	// ping open db to verify the connection has been established.
	err = p.db.Ping()
	if err != nil {
		panic(err)
	}
}
