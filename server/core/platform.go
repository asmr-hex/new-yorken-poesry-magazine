package core

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/env"
	_ "github.com/lib/pq"
)

// struct for storing in-memory statefulnessnessnesssnes for server
type Platform struct {
	*Logger
	Api    *API
	config *env.Config
	db     *sql.DB
}

func NewPlatform() *Platform {
	p := &Platform{
		Logger: NewLogger(os.Stdout),
		config: env.NewConfig(),
	}

	// connect to all the lovely things we must connect to in my life
	p.Connect()

	// setup db state, etc.
	p.Setup()

	// construct API and pass it the db connection handle set within Connect ---^
	p.Api = NewAPI(p.config, p.db)

	// print out server configuration
	if p.config.DevEnv {
		p.Info("Server running in Development mode")
	} else {
		p.Info("Server running in Production mode")
	}

	return p
}

// this is where it all happens.
// this is where we make the meaningful connections which will last for forever...
// just kidding, we will make more meaningful connections, don't you worry ( ⌒o⌒)人(⌒-⌒ )v
func (p *Platform) Connect() {
	var (
		err    error
		dbInfo string
	)

	// construct database info string required for connection
	dbInfo = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		p.config.DB.Host,
		p.config.DB.Username,
		p.config.DB.Password,
		p.config.DB.Name,
		"disable",
	)

	// open a connection to the database
	p.db, err = sql.Open(env.DB_DRIVER, dbInfo)
	if err != nil {
		panic(err)
	}

	// ping open db to verify the connection has been established.
	// otherwise (╥﹏╥)
	err = p.db.Ping()
	if err != nil {
		panic(err)
	}

	p.Success("Successful Connection -> %s", p.config.DB.Host)

	// if we connect to more services, we will do it below...
}

func (p *Platform) Setup() {
	// check to see if all tables have been created
	// for table := range DB_TABLE_NAMES {
	// 	//
	// }
}

func (p *Platform) Start() {
	// listen on quad-zero route with specified port yo (wait is this garbage?)
	addr := fmt.Sprintf("0.0.0.0:%s", p.config.Port)

	// here we gooooooooooo ʕつ•ᴥ•ʔつ
	err := http.ListenAndServe(addr, p.Api.Router)
	if err != nil {
		log.Fatal(err)
	}
}
