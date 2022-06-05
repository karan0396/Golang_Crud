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


	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func bgTask() {
	ticker := time.NewTicker(60 * time.Second)
	for range ticker.C {	
		 repo.DeleteUserbynTimes()
	}
}

var repo model.Repository

func main() {
	//Loading conifig file
	con, err := config.Load()
	if err != nil {
		logger.Logger.DPanic("config file is not load", zap.Error(err))
		return
	}
	//Calling repository for database active
	repo = model.NewRepository(con)

	//Initialize logger
	logger.IntializeLogger()
	r := mux.NewRouter()
	route.Routehandler(r, service.NewServ(repo))

	//Go routine for hard delete
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

