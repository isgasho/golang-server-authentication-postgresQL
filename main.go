package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
   "./signup"
   "./database"
	"github.com/gorilla/mux"
	"./token"
	"github.com/gorilla/context"
)

type message struct {
	Message string `json:"message"`
}
var db *sql.DB

type userId int

const secret string = "dasgasdgsd"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", token.AuthMiddleware(defaultHandler)).Methods("GET")
	r.HandleFunc("/signup", signUpHandler).Methods("POST")
	r.HandleFunc("/users", getUserHandler).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("POST")

	d, err := database.InitialDatabase("postgres://ogsbeoli:iUZvt7Teld42B8vAlDGTRjzdAVu7fZF9@isilo.db.elephantsql.com:5432/ogsbeoli")
	if err != nil{
		return
	}
	db = d
	log.Fatal(http.ListenAndServe(":7000", r))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoded := context.Get(r, "decoded")
	fmt.Println(decoded)
	w.WriteHeader(http.StatusOK)
	m := message{
		Message: "just health check",
	}
	json.NewEncoder(w).Encode(m)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	signupForm := signup.Signup{}
	var uID userId
	err := json.NewDecoder(r.Body).Decode(&signupForm)
	if err != nil {
		http.Error(w, "was not able to parse request body", http.StatusBadRequest)
		return
	}
	//hash password
	s, err := signupForm.HashPassword()
	if err != nil {
		http.Error(w, "was not able to hash your password", http.StatusBadRequest)
		return
	}

	err = db.QueryRow("INSERT INTO USERS (EMAIL, PASSWORD) VALUES ($1,$2) RETURNING id;", signupForm.Email,
		s).Scan(&uID)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Was not save your data to the database", http.StatusBadRequest)
		return
	}

	tokenString, err := token.CreateToken(signupForm.Email, s)

	if err != nil {
		http.Error(w, "was not able to generate a token for you fucker", http.StatusBadRequest)
		return
	}

	m := message{
		Message: fmt.Sprintf("it is saved, your token is %v and your id is %v", tokenString, uID),
	}

	json.NewEncoder(w).Encode(m)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	var s []signup.Signup
	c := make(chan signup.Signup)
	rows, err := db.Query("SELECT * FROM USERS")
	if err != nil {
		http.Error(w, "was not able to fetch data from the database", http.StatusBadRequest)
		return
	}

	go func() {
		var email string
		var password string
		var id int
		defer close(c)
		for rows.Next() {
			err = rows.Scan(&id, &email, &password)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "was not able to scan it", http.StatusBadRequest)
				return
			}
			sForm := signup.Signup{
				Email:    email,
				Password: password,
			}
			c <- sForm
		}
	}()

	for item := range c {
		s = append(s, item)
	}

	json.NewEncoder(w).Encode(s)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	s := signup.Signup{}
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "body was not included fucker", http.StatusBadRequest)
	}
	rows := db.QueryRow("SELECT * FROM USERS WHERE id=$1", 5)
	sD := signup.Signup{}
	var ID int
	err = rows.Scan(&ID, &sD.Email, &sD.Password)
	if err != nil {
		http.Error(w, "Was not able to scan it", http.StatusBadRequest)
		return
	}
	ok, err := s.ComparePassword(sD.Password, s.Password)
	if err != nil {
		http.Error(w, "password does not match", http.StatusBadRequest)
		return
	}
	if ok {
		tokenString, err := token.CreateToken(s.Email, sD.Password)
		if err != nil {
			http.Error(w, "was not able to generate a token", http.StatusBadRequest)
			return
		}

		m := message{
			Message: fmt.Sprintf("you are logged in, your token is %v", tokenString),
		}

		json.NewEncoder(w).Encode(m)

	}
}



