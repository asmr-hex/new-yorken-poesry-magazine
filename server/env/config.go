package env

import (
	"github.com/caarlos0/env"
)

const (
	DB_DRIVER = "postgres"
)

// server conf
type Config struct {
	DevEnv      bool   `env:"DEV_ENV" envDefault:"false"`
	Port        string `env:"SERVER_PORT" envDefault:"8080"`
	Magazine    MagazineConfig
	ExecContext ExecContext
	DB          DB
}

// tuneable parameters for the magazine...
type MagazineConfig struct {
	CommitteeSize          int     `env:"COMMITTEE_SIZE" envDefault:"10"`             // number of judges on a committee
	OpenSlotsPerIssue      int     `env:"OPEN_SLOTS_PER_ISSUE" envDefault:"13"`       // number of poems within a given issue
	CommitteeTurnoverRatio float64 `env:"COMMITTEE_TURNOVER_RATIO" envDefault:"0.50"` // percent of committee which is *replaced* on each issue cycle
	Randomness             float64 `env:"RANDOMNESS" envDefault:"0.10"`               // uhhh, randomness of parameters.....
	MetaRandomness         float64 `env:"META_RANDOMNESS" envDefault:"0.05"`          // uhhh, randomness of the randomness parameter
	Pretension             float64 `env:"PRETENSION" envDefault:"0.5"`                // percent of new judges who are chosen purely on the "quality/quantity" of their published work (i.e. how much to *not* select for "underdogs")
}

type ExecContext struct {
	Dir      string `env:"EXEC_DIR" envDefault:"/tmp"`
	MountDir string `env:"EXEC_MOUNT_DIR" envDefault:"/tmp"`
}

// DB conf
type DB struct {
	Host     string `env:"DB_HOST" envDefault:"postgres"`
	Username string `env:"POSTGRES_USER" envDefault:"wintermute"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"t0b30rn0tt0b3"`
	Name     string `env:"POSTGRES_DB" envDefault:"nypm"`
}

// parse environment variables into Config
func NewConfig() *Config {
	conf := &Config{}

	// parse server conf
	err := env.Parse(conf)
	if err != nil {
		panic(err)
	}

	// parse Magazine conf
	err = env.Parse(&conf.Magazine)
	if err != nil {
		panic(err)
	}

	// parse DB conf
	err = env.Parse(&conf.DB)
	if err != nil {
		panic(err)
	}

	// parse execution context
	err = env.Parse(&conf.ExecContext)
	if err != nil {
		panic(err)
	}

	// parse other confs (when they present themselves...?)
	// in the meantime, let's take a walk or, perhaps be emotional with one another
	// (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥)

	return conf
}

/*

   testing configs

*/

// test server conf
type TestConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
	DB   TestDB
}

// DB conf
type TestDB struct {
	Host     string `env:"TEST_DB_HOST" envDefault:"postgres"`
	Username string `env:"TEST_POSTGRES_USER" envDefault:"wintermute"`
	Password string `env:"TEST_POSTGRES_PASSWORD" envDefault:"t0b30rn0tt0b3"`
	Name     string `env:"TEST_POSTGRES_DB" envDefault:"nypm_test"`
}

// parse environment variables into Config
func NewTestConfig() *TestConfig {
	conf := &TestConfig{}

	// parse server conf
	err := env.Parse(conf)
	if err != nil {
		panic(err)
	}

	// parse DB conf
	err = env.Parse(&conf.DB)
	if err != nil {
		panic(err)
	}

	return conf
}
