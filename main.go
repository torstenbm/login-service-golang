package main

import (
	"fmt"
	"log"
	"login-service/controller"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.POST("/register", controller.RegisterUser)
	router.POST("/login", controller.LoginUser)
	router.POST("/logout", controller.LogoutUser)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
