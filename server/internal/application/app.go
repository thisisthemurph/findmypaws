package application

import (
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"paws/internal/routes"
	"paws/internal/store"

	_ "github.com/lib/pq"
)

type App struct {
	Store  *store.PostgresStore
	Router *routes.Router
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

	return &App{
		Logger: logger,
		Config: config,
	}, nil
}

func (app *App) Build() error {
	app.Logger.Info("building app")

	if err := app.configureStores(); err != nil {
		return err
	}

	r := routes.NewRouter(app.Store, app.Logger)
	app.Router = r
	return nil
}

func (app *App) configureStores() error {
	app.Logger.Info("configuring stores")

	db, err := sqlx.Open("postgres", app.Config.Database.ConnectionString)
	if err != nil {
		return err
	}

	s := store.NewPostgresStore(db)
	app.Store = s
	return nil
}
