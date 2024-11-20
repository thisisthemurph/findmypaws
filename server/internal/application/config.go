package application

import "fmt"

type DatabaseConfig struct {
	ConnectionString string
}

type ClerkConfig struct {
	Secret        string
	SigningSecret string
}

type AppConfig struct {
	Host          string
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

	return AppConfig{
		Host:          get("HOST"),
		ClientBaseURL: get("CLIENT_BASE_URL"),
		Clerk: ClerkConfig{
			Secret:        get("CLERK_SECRET_KEY"),
			SigningSecret: get("CLERK_SIGNING_SECRET"),
		},
		Database: DatabaseConfig{
			ConnectionString: get("DATABASE_CONNECTION_STRING"),
		},
	}
}
