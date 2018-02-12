package core

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	_ "github.com/lib/pq"
)

const (
	DB_DRIVER = "postgres"
)

// struct for storing in-memory statefulnessnessnesssnes for server
type Platform struct {
	Api    *API
	config *Config
	db     *sql.DB
	logger *log.Logger
}

func NewPlatform() *Platform {
	p := &Platform{
		config: NewConfig(),
		logger: log.New(os.Stdout, "", 0), // "you *really* don't know what this zero means?" -- ex-coworker
	}

	// connect to all the lovely things we must connect to in my life
	p.Connect()

	// construct API and pass it the db connection handle set within Connect ---^
	p.Api = NewAPI(p.db)

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
	p.db, err = sql.Open(DB_DRIVER, dbInfo)
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

// some Platform utils functions
// i know what your saying right now, YAGNI, but i can't help myself (‾⌣‾)♉
// i could probably write all these in a better way by implementing a custom io.Writer
// for each case (success, info, error) and wrap the text in a color of choice and specify
// a prefix on an individual basis, but i don't feel like it. But i probably could have
// done that during the time it took me to write this. #designdecisions
func (p *Platform) Success(format string, v ...interface{}) {
	p.logger.Print(color.GreenString(fmt.Sprintf(format, v...)))
}

func (p *Platform) Info(format string, v ...interface{}) {
	p.logger.Print(color.BlueString(fmt.Sprintf(format, v...)))
}

func (p *Platform) Error(format string, v ...interface{}) {
	p.logger.Print(color.RedString(fmt.Sprintf(format, v...)))
}
