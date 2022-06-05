package middleware

import (
	"api/pkg/token"
	"api/pkg/logger"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go"
)
func init() {
	logger.IntializeLogger()
	
}
//Auhtorization the login process
func Authorization(s http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ss := r.Header.Get("Authorization")
		if ss == "" {
			logger.Logger.Info("Checking header")
			w.WriteHeader(http.StatusBadRequest)
			logger.Logger.DPanic("not Authorized", zap.String("header", ss))
			return
		}
		arr := strings.Fields(ss)
		if len(arr) != 2 {
			logger.Logger.Info("checking array")
			w.WriteHeader(http.StatusBadRequest)
			logger.Logger.DPanic("not Authorized", zap.String("header", ss))
			return
		}

		if strings.ToLower(arr[0]) != "bearer" {
			logger.Logger.Info("checking lower case of bearer token")
			w.WriteHeader(http.StatusBadRequest)
			logger.Logger.DPanic("not Authorized", zap.String("header", ss))
			return
		}
		tkstr := arr[1]

		tkn, err := token.ParseToken(tkstr)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				logger.Logger.Info("Signature is invalid")
				logger.Logger.DPanic("signture is not valid", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			logger.Logger.Info("Token is not valid")
			logger.Logger.DPanic("token is not vaid", zap.Bool("token is valid", tkn.Valid))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		s(w, r)
	}
}
