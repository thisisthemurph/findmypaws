package application

import (
	"fmt"
	"strconv"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

func (e Environment) String() string {
	return string(e)
}

func (e Environment) IsDevelopment() bool {
	return e == Development
}

type DatabaseConfig struct {
	Name             string
	ConnectionString string
	ForceMigration   bool
}

type ClerkConfig struct {
	Secret        string
	SigningSecret string
}

type AppConfig struct {
	Host          string
	Environment   Environment
	ClientBaseURL string
	Database      DatabaseConfig
	Clerk         ClerkConfig
}

func NewAppConfig(getFunc func(string) string) AppConfig {
	get := func(k string) string {
		v := getFunc(k)
		if v == "" {
			panic(fmt.Sprintf("Evironment variable %q not found", k))
		}
		return v
	}

	forceDatabaseMigration, err := strconv.ParseBool(get("DATABASE_FORCE_MIGRATION"))
	if err != nil {
		panic(err)
	}

	return AppConfig{
		Host:          get("HOST"),
		Environment:   Environment(get("ENVIRONMENT")),
		ClientBaseURL: get("CLIENT_BASE_URL"),
		Clerk: ClerkConfig{
			Secret:        get("CLERK_SECRET_KEY"),
			SigningSecret: get("CLERK_SIGNING_SECRET"),
		},
		Database: DatabaseConfig{
			Name:             get("DATABASE_NAME"),
			ConnectionString: get("DATABASE_CONNECTION_STRING"),
			ForceMigration:   forceDatabaseMigration,
		},
	}
}
