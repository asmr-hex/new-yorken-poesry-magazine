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
	Api    *API
	config *Config
	db     *sql.DB
}

func NewPlatform() *Platform {
	p := &Platform{
		Api: NewAPI(),
		config: &Config{
			User:     "wintermute",
			Password: "t0b30rn0tt0b3",
			DBName:   "nypm",
		},
	}

	// TODO retry until db connects?
	p.Connect()

	return p
}

func (p *Platform) Connect() {
	var (
		err    error
		dbInfo string
	)

	// TODO parameterize "host=db" so we have a single source of truth config or something
	// actualy load from .env
	dbInfo = fmt.Sprintf("host=postgres user=%s password=%s dbname=%s sslmode=%s",
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

	fmt.Println("CONNECTED TO DB YOO")
}
