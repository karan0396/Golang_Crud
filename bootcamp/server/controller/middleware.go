package controller

import (
	"bootcamp/model"
	// "go/token"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)


func Authorize(s http.HandlerFunc)http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		c,err:=r.Cookie("token")
		if err!=nil{
			if err == http.ErrNoCookie{
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return

		}

		tkstr:=c.Value

		claims:= &model.Claims{}

		tkn,err:=jwt.ParseWithClaims(tkstr,claims,func (token *jwt.Token)(interface{},error) {
			return jwtKey,nil
		})

		if err!=nil{
			if err==jwt.ErrSignatureInvalid{
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid{
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		s(w,r)
	}
}




// func Cors(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Set("Access-Control-Allow-Origin", "*")                                                            
// 			w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") 
// 			// w.Header().Add("Access-Control-Allow-Credentials", "true")                                                   
// 			w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                             
// 			w.Header().Set("content-type", "application/json;charset=UTF-8")
// 			if r.Method == "OPTIONS" {
// 					w.WriteHeader(http.StatusNoContent)
// 					return
// 			}
// 			next.ServeHTTP(w, r)
// 	})
// }