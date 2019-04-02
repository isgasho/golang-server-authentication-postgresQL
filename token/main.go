package token

import (
	"github.com/dgrijalva/jwt-go"
	"strings"
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/gorilla/context"
)

type message struct{
	Message string `json:"string"`
}

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
//AuthMiddleware for authentication
func AuthMiddleware(next http.HandlerFunc)http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get header
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != ""{
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{},error){
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("There was an error")
                    }
                    return []byte(secret), nil
				})

				if err != nil {
					m := message{"Invalid authrozation token"}
                    json.NewEncoder(w).Encode(m)
                    return
				}
				
				if token.Valid {
                    context.Set(r, "decoded", token.Claims)
                    next(w, r)
                } else {
					m := message{"Invalid authrozation token"}
                    json.NewEncoder(w).Encode(m)
                }
			}
		}else {
			m := message{"authorized header is needed"}
			json.NewEncoder(w).Encode(m)
        }
	})
}