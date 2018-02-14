package env

import (
	"github.com/caarlos0/env"
)

const (
	DB_DRIVER = "postgres"
)

// server conf
type Config struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
	DB   DB
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

	// parse DB conf
	err = env.Parse(&conf.DB)
	if err != nil {
		panic(err)
	}

	// parse other confs (when they present themselves...?)
	// in the meantime, let's take a walk or, perhaps be emotional with one another
	// (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥) (♥_♥)

	return conf
}
