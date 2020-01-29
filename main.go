package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/signup", signup).Methods("POST")

	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/protected", TokenVerifyMiddleware(ProtectedEndPoint)).Methods("POST")
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func signup(w http.ResponseWriter, r *http.Request) {

	var user User
	var error Error
	json.NewDecoder(r.Body).Decode(&user)

	spew.Dump(user)

	if user.Email == "" {
		error.Message = "Email is missing"

		//send bad request
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is missing"

		//send bad request
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}
	fmt.Println("Printing error")
	hash, err:=bcrypt.GenerateFromPassword([]byte(user.Password),10)
	fmt.Println(err)
	if err!=nil{
		log.Fatal("error")
	}
	fmt.Println(hash)
	user.Password= string(hash) 
	fmt.Println("user password ",user.Password)
}
func login(w http.ResponseWriter, r *http.Request) {

	var user User
	 
	json.NewDecoder(r.Body).Decode(&user)
	data:= GenerateToken(user)
	spew.Dump(data)
	//w.Write([]byte("runnning login"))

}
func ProtectedEndPoint(w http.ResponseWriter, r *http.Request) {

}
func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return nil
}
func GenerateToken(user User)(string){

   secret :="secret"

 token:= jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
	  "email":user.Email,
	   "iss":"course",
  })

  tokenString,err := token.SignedString([]byte(secret))

		if err !=nil {
			log.Fatal(err)
		}
	 
   		return tokenString
}
