package main

import (
	"net/http"
	"paws/internal/application"
)

func main() {
	app, err := application.NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Build(); err != nil {
		panic(err)
	}

	http.ListenAndServe(app.Config.Host, app.ServerMux)
}
