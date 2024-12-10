package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"paws/internal/application"
	"paws/internal/routes"
	"paws/pkg/migrator"
)

func main() {
	logger := log.New(os.Stdout, "API: ", log.LstdFlags|log.Lshortfile)
	if err := run(logger); err != nil {
		logger.Fatal(err)
	}
}

func run(logger *log.Logger) error {
	logger.Println("Starting API server...")
	app, err := application.NewApp()
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	if err := migrateDatabase(app); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	logger.Println("Building application...")
	if err := app.Build(); err != nil {
		return fmt.Errorf("failed to build application: %w", err)
	}

	logger.Println("Setting up routes...")
	mux := routes.BuildRoutesServerMux(app)

	logger.Printf("Starting server at %s...", app.Config.Host)
	if err := http.ListenAndServe(app.Config.Host, mux); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func migrateDatabase(app *application.App) error {
	dbConfig := app.Config.Database
	environment := app.Config.Environment

	if environment.IsDevelopment() && !dbConfig.ForceMigration {
		app.Logger.Warn("skipping database migration", "environment", environment)
		return nil
	}

	app.Logger.Info("Migrating database", "environment", environment, "force", dbConfig.ForceMigration)
	db, err := sql.Open("postgres", dbConfig.ConnectionString)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}
	defer db.Close()

	m := migrator.NewPostgresMigrator(db, dbConfig.Name, migrator.DefaultMigrationPath).WithLogger(app.Logger)
	return m.Migrate(migrator.MigrationDirectionUp)
}
