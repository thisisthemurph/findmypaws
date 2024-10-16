package application

type DatabaseConfig struct {
	ConnectionString string
}

type SupabaseConfig struct {
	URL            string
	AnonKey        string
	ServiceRoleKey string
	JWTSecret      string
}

type AppConfig struct {
	Host          string
	ClientBaseURL string
	Database      DatabaseConfig
	Supabase      SupabaseConfig
}

func NewAppConfig(get func(string) string) AppConfig {
	return AppConfig{
		Host:          get("HOST"),
		ClientBaseURL: get("CLIENT_BASE_URL"),
		Database: DatabaseConfig{
			ConnectionString: get("DATABASE_CONNECTION_STRING"),
		},
		Supabase: SupabaseConfig{
			URL:            get("SUPABASE_URL"),
			AnonKey:        get("SUPABASE_ANON_KEY"),
			ServiceRoleKey: get("SUPABASE_SERVICE_ROLE_KEY"),
			JWTSecret:      get("SUPABASE_JWT_SECRET"),
		},
	}
}
