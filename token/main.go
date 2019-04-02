package token

import (
	"github.com/dgrijalva/jwt-go"
)

const secret string = "What is up brother"

//CreateToken is a util function
func CreateToken(email,password string) (string,error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    email,
		"password": password,
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}