package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"paws/internal/application"
	"paws/internal/routes"
	"paws/pkg/migrator"
)

func main() {
	app, err := application.NewApp()
	if err != nil {
		panic(err)
	}

	if err := migrateDatabase(app); err != nil {
		log.Fatal("failed to migrate database", err)
	}

	if err := app.Build(); err != nil {
		log.Fatal(err)
	}

	mux := routes.BuildRoutesServerMux(app)
	http.ListenAndServe(app.Config.Host, mux)
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
