package main

import (
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"jwt-course/controllers"
	"jwt-course/driver"
	"jwt-course/utils"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "main: ", log.LstdFlags)
	err := gotenv.Load()
	if err != nil {
		logger.Fatalln(err)
	}
}

func main() {
	db := driver.ConnectDB()

	router := mux.NewRouter()
	controller := controllers.NewController(db)
	router.HandleFunc("/signup", controller.SignUp()).Methods("POST")
	router.HandleFunc("/login", controller.Login()).Methods("POST")
	router.HandleFunc("/protected", utils.TokenVerifyMiddleware(controller.Protected())).Methods("GET")

	logger.Println("starting server on port 8000")

	log.Fatalln(http.ListenAndServe(":8000", router))
}
