package controller

import (
	"encoding/json"
	"login-service/model"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func writeResponse(w http.ResponseWriter, value string, message string) {
	response := make(map[string]string)
	response[value] = message

	parsedResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Bad request", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(parsedResponse)
	return
}

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

		taken, err := repository.IsUserNameTaken(un)

		if taken {
			writeResponse(w, "error", "Username is taken: "+err.Error())
			return
		}
		uf := model.UserFactory{}
		u := uf.CreateUser(un, f, l, p, time.Now().Local())

		err = repository.WriteUserToDb(u)
		if err != nil {
			writeResponse(w, "error", "There was en error writing user to db: "+err.Error())
			return
		}

		// Sending access response
		writeResponse(w, "access", "granted") // Should switch for security token when we learn about security from Carlos
		return
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

		user, err := repository.GetUserFromDb(un)
		if err != nil {
			writeResponse(w, "error", "There was en error retrieving user from db: "+err.Error())
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(p))
		if err != nil {
			writeResponse(w, "error", "Username and password does not match: "+err.Error())
			return
		}

		// Sending access response
		writeResponse(w, "access", "granted") // Should switch for security token when we learn about security from Carlos
		return
	}
}

// LogoutUser : Handles logout. As of now this is handled by front-end, might change in the future if I learn that it's good
func LogoutUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
