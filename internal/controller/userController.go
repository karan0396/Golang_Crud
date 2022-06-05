package controller

import (
	"api/internal/model"
	"api/internal/service"
	"api/pkg/logger"
	"api/pkg/parse"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"


	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type register struct {
	svc service.NewService
}

func init() {
	logger.IntializeLogger()
}

func NewRegister(svc service.NewService) *register {
	reg := &register{svc}
	return reg
}




func (reg *register) Login(w http.ResponseWriter, r *http.Request) {
	var data model.Credential

	json.NewDecoder(r.Body).Decode(&data)
	tokenString, error := reg.svc.LoginService(data)
	if error != nil {
		http.Error(w, "email or password is incorect", http.StatusNotFound)
		return
	}

	e := service.Signinresponse{
		Token:   tokenString,
		Message: "Login Success",
	}

	json.NewEncoder(w).Encode(e)
}

func (reg *register) CreateUser(w http.ResponseWriter, r *http.Request) {
	var data model.User
	err := parse.Parse(w, r, &data)
	if err != nil {
		logger.Logger.DPanic("parse is not done", zap.Error(err))
		http.Error(w,"Data is not added Correctly",http.StatusNotAcceptable)
		return
	}
	// fmt.Println(data)

	err = reg.svc.CreateService(data)
	if err != nil {
		http.Error(w, "data not added", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Data Added")
}

func (reg *register) GetUser(w http.ResponseWriter, r *http.Request) {
	p := model.Params{
		Archieved: r.URL.Query().Get("archived"),
		Id:        r.URL.Query().Get("id"),
		Name:      r.URL.Query().Get("name"),
		Email:     r.URL.Query().Get("email"),
		Sort:      r.URL.Query().Get("sort"),
		Order:     r.URL.Query().Get("order"),
		Page:      r.URL.Query().Get("page"),
		Limit:     r.URL.Query().Get("limit"),
	}

	newuser, err := reg.svc.GetService(p)
	if err != nil {
		logger.Logger.DPanic("user is not created", zap.Error(err))
		http.Error(w,"Data is not created",http.StatusBadRequest)
	}
	res, _ := json.Marshal(newuser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (reg *register) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	Id, err := strconv.Atoi(userId)
	if err != nil {
		logger.Logger.Debug("Don't converted into int", zap.Error(err))
		http.Error(w,"Error in getting Id",http.StatusInternalServerError)
	}
	err = reg.svc.DeleteService(Id)
	if err != nil {
		logger.Logger.DPanic("Not deleted", zap.Error(err))
		http.Error(w,"Not Deleted",http.StatusBadGateway)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}

func (reg *register) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var data model.User
	err := parse.Parse(w, r, &data)
	if err != nil {
		logger.Logger.DPanic("parse is not done", zap.Error(err))
		http.Error(w,"parsing &validation is not done",http.StatusNotAcceptable)
	}
	fmt.Println(data)

	vars := mux.Vars(r)
	userId := vars["id"]
	Id, err := strconv.Atoi(userId)
	if err != nil {
		logger.Logger.Debug("Don't converted into int", zap.Error(err))
		http.Error(w,"id not converted to int",http.StatusInternalServerError)
	}
	data.ID = Id
	err = reg.svc.UpdateService(data)
	if err != nil {
		logger.Logger.DPanic("not Updated", zap.Error(err))
		http.Error(w,"Not Updated",http.StatusBadGateway)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Data Updated")
}
