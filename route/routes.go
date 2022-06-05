package route

import (
	"api/internal/controller"
	"api/internal/middleware"
	"api/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)


func Routehandler(router *mux.Router,svc service.NewService){
	reg := controller.NewRegister(svc)
	router.HandleFunc("/login", reg.Login).Methods("Post")
	router.HandleFunc("/user", middleware.Authorization(reg.CreateUser)).Methods("POST")
	router.HandleFunc("/user", reg.GetUser).Methods("Get")
	router.HandleFunc("/user/{id}", middleware.Authorization(reg.DeleteUser)).Methods("Delete")
	router.HandleFunc("/user/{id}", middleware.Authorization(reg.UpdateUser)).Methods("Patch")
	router.Handle("/favicon.ico", http.NotFoundHandler())
}