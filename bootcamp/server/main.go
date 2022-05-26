package main

import (
	// "bootcamp/controller"

	"bootcamp/route"
	"net/http"

	"github.com/rs/cors"
)

func main() {

	r := route.Route()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost},
		AllowedHeaders: []string{"Authorization","Content-Type"},
		AllowCredentials: true,
		

	})

	handler := c.Handler(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	srv.ListenAndServe()

}
