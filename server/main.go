package main

import (
	"api/config"
	"api/internal/model"
	"api/internal/service"
	"api/pkg/logger"
	"api/route"
	"fmt"
	"net/http"
	"time"

	// "time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
	// "go.uber.org/zap"
)

func bgTask() {
	ticker := time.NewTicker(3 * time.Hour)
	for range ticker.C {		
		 repo.DeleteUserbynTimes()
	}
}

var repo model.Repository

func main() {
	con, err := config.Load()
	if err != nil {
		logger.Logger.DPanic("config file is not load", zap.Error(err))
		return
	}

	repo = model.NewRepository(con)


	logger.IntializeLogger()
	r := mux.NewRouter()
	route.Routehandler(r, service.NewServ(repo))
	
	go bgTask()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%v", con.Server.Port),
		Handler: handler,
	}

	srv.ListenAndServe()
	select {}
}

