package token

import (
	"api/pkg/logger"
	// "encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	// "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("Suck_it")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
//Creating Token
func CreateToken(email string, expiringTime time.Time) (string, error) {
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiringTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	logger.Logger.Info("Creating Token")
	if err != nil {
		logger.Logger.DPanic("key is not done", zap.Error(err))
		return "", err
	}
	return tokenString, nil
}
// parsing token
func ParseToken(token string) (*jwt.Token, error) {
	claims := Claims{}
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		logger.Logger.Info("Parsing Token")
		return jwtKey, nil
	})
	if err != nil {
		logger.Logger.DPanic("rerieving token is not done", zap.Error(err))
		return nil, err
	}
	return tkn, nil

}