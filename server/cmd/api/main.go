package main

import (
	"log"
	"net/http"
	"paws/internal/application"
	"paws/internal/routes"
)

func main() {
	app, err := application.NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Build(); err != nil {
		log.Fatal(err)
	}

	mux := routes.BuildRoutesServerMux(app)
	http.ListenAndServe(app.Config.Host, mux)
}
