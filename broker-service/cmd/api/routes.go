//This is all of the routes for the application
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (appli *Config) routes() http.Handler {
	/*multiplexer HandleFunction*/
	mux := chi.NewRouter()

	//specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accpet", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, //allow to send request to deal with credentials
		MaxAge:           300,  //300 is a good default
	}))

	//check its heartbeat to make sure whether this service is alive
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/", appli.brokerHandler)

	//real route to listen to grpc
	mux.Post("/log-grpc", appli.LogViaGRPC)

	mux.Post("/handle", appli.HandleSubmission)

	return mux
}
