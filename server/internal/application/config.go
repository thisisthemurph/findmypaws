package application

type DatabaseConfig struct {
	ConnectionString string
}

type ClerkConfig struct {
	Secret string
}

type AppConfig struct {
	Host          string
	ClientBaseURL string
	Database      DatabaseConfig
	Clerk         ClerkConfig
}

func NewAppConfig(get func(string) string) AppConfig {
	return AppConfig{
		Host:          get("HOST"),
		ClientBaseURL: get("CLIENT_BASE_URL"),
		Clerk: ClerkConfig{
			Secret: get("CLERK_SECRET_KEY"),
		},
		Database: DatabaseConfig{
			ConnectionString: get("DATABASE_CONNECTION_STRING"),
		},
	}
}
