package route

import (
	"api/internal/controller"
	"api/internal/middleware"
	"api/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

//Routes
func Routehandler(router *mux.Router,svc service.NewService){
	reg := controller.NewRegister(svc)
	router.HandleFunc("/login", reg.Login).Methods("Post")   									//for login
	router.HandleFunc("/user", middleware.Authorization(reg.CreateUser)).Methods("POST")  		//for Creating User
	router.HandleFunc("/user", reg.GetUser).Methods("Get") 										//for getiing all user
	router.HandleFunc("/user/{id}", middleware.Authorization(reg.DeleteUser)).Methods("Delete")	//for deleteing user
	router.HandleFunc("/user/{id}", middleware.Authorization(reg.UpdateUser)).Methods("Patch")	//for updating 
	router.Handle("/favicon.ico", http.NotFoundHandler())    
}