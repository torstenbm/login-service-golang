package controller

import (
	"encoding/json"
	"login-service2/model"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser : Handles registering of users, writes user to DB and sends back response
func RegisterUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	if r.Method == http.MethodPost {
		repository := model.UserRepository{}

		un := r.FormValue("username")
		f := r.FormValue("firstname")
		l := r.FormValue("lastname")
		p := r.FormValue("password")

		taken := repository.IsUserNameTaken(un)

		if taken {
			http.Error(w, "Username is taken", http.StatusForbidden)
			return
		}
		uf := model.UserFactory{}
		u := uf.CreateUser(un, f, l, p, time.Now().Local())

		repository.WriteUserToDb(u)

		// Sending response

		response := make(map[string]string) // Should switch for security token when we learn about security from Carlos
		response["access"] = "granted"
		parsedResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Bad request", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(parsedResponse)
	}

}

// LoginUser : Handles login
func LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	repository := model.UserRepository{}

	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		p := r.FormValue("password")

		user := repository.GetUserFromDb(un)

		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(p))
		if err != nil {
			http.Error(w, "Username and password does not match", http.StatusForbidden)
			return
		}

		// Sending response
		response := make(map[string]string) // Should switch for security token when we learn about security from Carlos
		response["access"] = "granted"

		parsedResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Bad request", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(parsedResponse)
	}
}

// LogoutUser : Handles logout. As of now this is handled by front-end, might change in the future if I learn that it's good
func LogoutUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
