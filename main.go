package main

import (
	"fmt"
	"log"
	"login-service/controller"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "This is the login microservice build in go\n")
}
func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/register", controller.RegisterUser)
	router.POST("/login", controller.LoginUser)
	router.POST("/logout", controller.LogoutUser)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
