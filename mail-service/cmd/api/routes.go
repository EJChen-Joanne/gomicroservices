package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (appli *Config) routes() http.Handler {
	//add a new router
	mux := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	mux.Use(cors.Handler)

	mux.Use(middleware.Heartbeat("/ping"))

	//add handler to receive the request for sending email
	mux.Post("/send", appli.SendMail)
	return mux
}
