package main

import (
	"fmt"
	"log"
	"login-service/controller"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "This is the login microservice build in go\n")
}
func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/register", controller.RegisterUser)
	router.POST("/login", controller.LoginUser)
	router.POST("/logout", controller.LogoutUser)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
