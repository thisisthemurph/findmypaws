package main

import "paws/internal/application"

func main() {
	app, err := application.NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Build(); err != nil {
		panic(err)
	}

	app.Router.Logger.Fatal(app.Router.Start(app.Config.Host))
}
