package application

type DatabaseConfig struct {
	ConnectionString string
}

type AppConfig struct {
	Host     string
	Database DatabaseConfig
}

func NewAppConfig(get func(string) string) AppConfig {
	return AppConfig{
		Host: get("HOST"),
		Database: DatabaseConfig{
			ConnectionString: get("DATABASE_CONNECTION_STRING"),
		},
	}
}
