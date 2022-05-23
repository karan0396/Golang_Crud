package controller

import (
	// "bootcamp/database"
	"bootcamp/config"
	"bootcamp/logger"
	"bootcamp/model"
	"database/sql"

	// "bootcamp/util"
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)
var db *sql.DB
var jwtKey = []byte("my_secret_key")


func Signin(w http.ResponseWriter,r *http.Request){
	db=config.GetDb()
	var cred model.Credential
	var password string   //storing password

	//get credential by requesting

	json.NewDecoder(r.Body).Decode(&cred)

	err:=db.QueryRow(`select password from user where email=?`,cred.Email).Scan(&password)
	if err!=nil{
		http.Error(w,"Email not Found",http.StatusNotFound)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//matching credential

	err= bcrypt.CompareHashAndPassword([]byte(password), []byte(cred.Password))
	if err!=nil{
		http.Error(w,"Password Not Match",http.StatusUnauthorized)
		return
	}

	//expiration time of cookie
	expirationTime := time.Now().Add(10 * time.Minute)


	//making struct for saving cookie
	claims:= &model.Claims{
		Email:cred.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// fmt.Println(claims)

	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	tokenString,err:=token.SignedString(jwtKey)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,&http.Cookie{
		Name:"token",
		Value: tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode("Login Succcess")
}


func Logout(w http.ResponseWriter,r *http.Request){
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Already logged out", http.StatusBadRequest)
		return
	}

	cookie = &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode("logged out")
}





func GetUser(w http.ResponseWriter, r *http.Request) {
	newuser:=model.GetAllUser(r)
	res, _ :=json.Marshal(newuser)
	
	w.Header().Set("Content-Type","pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}


func DeleteUser(w http.ResponseWriter,r *http.Request){
	vars:=mux.Vars(r)
	userId:=vars["id"]
	Id, err :=  strconv.ParseInt(userId,0,0)
	if err!=nil{
		logger.ErrorLogger.Println("error while parsing")
	}
	model.DeleteUser(Id)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}



func CreatUser(w http.ResponseWriter,r *http.Request){
	var body model.ReqUser
	x := r.Body
	err:=json.NewDecoder(x).Decode(&body)
	if err!=nil{
		http.Error(w,"this is not encoding",http.StatusInternalServerError)
		return
	}
	if !ValidationReq(&body){
		http.Error(w,"Please check the entered data",http.StatusBadRequest)
		return
	}
	body.CreatUser()
	w.Header().Set("Content-Type","pkglication/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Data Added")
}

func UpdateUser(w http.ResponseWriter,r *http.Request){
	updated:=time.Now()
	var body model.ReqUser
	err:=json.NewDecoder(r.Body).Decode(&body)
	if err!=nil{
		http.Error(w,"update time decoding",http.StatusInternalServerError)
		return
	}
	if !ValidationReq(&body){
		http.Error(w,"Please check the entered data",http.StatusBadRequest)
		return
	}

	vars:=mux.Vars(r)
	userId:=vars["id"]
	Id, err :=  strconv.Atoi(userId)
	if err!=nil{
		logger.ErrorLogger.Println("error while parsing")
	}

	var email string
	err = db.QueryRow("SELECT email FROM users where user_id=?", Id).Scan(&email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.ErrorLogger.Println(err)
		return
	}

	valid := false
	if len(body.FName)+len(body.Lname) < 30 && len(body.Email) < 20 {
		valid = true
	}
	if valid==false{
		http.Error(w,"Enter data of valid length",http.StatusBadRequest)
		logger.ErrorLogger.Println("Email or password exceeds maximum length")
		return
	}

	//check if new email is already in use 
	if email != body.Email && body.Email != "" {
		var countEmail int
		db.QueryRow("SELECT COUNT(email) FROM users where email=?", body.Email).Scan(&countEmail)
		if countEmail != 0 {
			http.Error(w, "Email already in use", http.StatusBadRequest)
			return
		}
	}


	// store retrive data in struct 
	userDetail,db:=model.GetUserByid(Id)
	if body.FName!=""{
		userDetail.FName=body.FName
	}
	if body.Lname!=""{
		userDetail.Lname=body.Lname
	}
	if body.Email!=""{
		userDetail.Email=body.Email
	}
	if body.Dob!=""{
		userDetail.Dob=body.Dob
	}

	//updating value
	_,err=db.Exec(`Update user Set firstname=?,lastname=?,email=?,dob=?,updated_at=? where id=?`,userDetail.FName,userDetail.Lname,userDetail.Email,userDetail.Dob,updated,Id)
	if err!=nil{
		logger.ErrorLogger.Println("update not work")
	}

	//sending response
	res,_:=json.Marshal(userDetail)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}