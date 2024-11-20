package application

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"paws/internal/chat"
	"paws/internal/repository"
)

type App struct {
	DB           *sqlx.DB
	ChatManager  *chat.Manager
	Repositories *repository.Repositories
	//ServerMux    *http.ServeMux
	Logger *slog.Logger
	Config AppConfig
}

func NewApp() (*App, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := NewAppConfig(os.Getenv)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	clerk.SetKey(config.Clerk.Secret)

	return &App{
		Logger: logger,
		Config: config,
	}, nil
}

func (app *App) Build() error {
	app.Logger.Info("building app")

	if err := app.configureDatabase(); err != nil {
		return err
	}
	app.configureRepositories()
	app.configureChatManager()

	//app.ServerMux = routes.BuildRoutesServerMux(app)
	return nil
}

func (app *App) configureDatabase() error {
	app.Logger.Info("configuring stores")

	db, err := sqlx.Open("postgres", app.Config.Database.ConnectionString)
	if err != nil {
		return fmt.Errorf("could not open to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	app.DB = db
	return nil
}

func (app *App) configureRepositories() {
	app.Logger.Info("configuring repositories")
	app.Repositories = repository.NewRepositories(app.DB)
}

func (app *App) configureChatManager() {
	app.Logger.Info("configuring chat manager")
	app.ChatManager = chat.NewManager(app.DB, app.Logger)
}
