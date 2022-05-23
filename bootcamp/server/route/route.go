package route

import (
	"bootcamp/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Route()*mux.Router{
	r:=mux.NewRouter()
	r.Use(controller.Cors)

	r.HandleFunc("/user",controller.Authorize(controller.GetUser)).Methods("Get")
	// r.HandleFunc("/user/{id}",GetUser).Methods("Get")
	r.HandleFunc("/user",controller.CreatUser).Methods("Post")
	r.HandleFunc("/user/{id}",controller.Authorize(controller.DeleteUser)).Methods("DELETE")
	r.HandleFunc("/user/{id}",controller.Authorize(controller.UpdateUser)).Methods("PATCH")
	r.HandleFunc("/signin",controller.Signin).Methods("Post")
	r.HandleFunc("/logout",controller.Logout).Methods("Post")
	log.Fatal(http.ListenAndServe(":8000",r))
	return r
}
