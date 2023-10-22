package env

import "fmt"

var Env Environment

type Environment struct {
	Postgres struct {
		User     string `env:"USER,default=postgres"`
		Password string `env:"PASSWORD,default=postgres"`
		DbName   string `env:"DB_NAME,default=postgresL0"`
	}

	Nats struct {

	}
}

func (e Environment) GetPostgresURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@localhost:5432/%s",
		e.Postgres.User,
		e.Postgres.Password,
		e.Postgres.DbName,
	)
}