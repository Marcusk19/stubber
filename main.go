package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/Marcusk19/stubber/controller"
	"github.com/Marcusk19/stubber/util"
)

type Movie struct {
	ID     string `json:"Id"`
	Title  string `json:"Title"`
	Desc   string `json:"desc"`
	Rating string `json:"rating"`
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	// myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/v1/movies", controller.MoviesByRating).Methods("GET")
	myRouter.HandleFunc("/api/v1/movies/{id}", controller.GetMovie).Methods("GET")
	myRouter.HandleFunc("/api/v1/movies", controller.InsertMovie).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/api/v1/movies", controller.UpdateMovie).Methods("POST")
	myRouter.HandleFunc("/api/v1/movies/delete/{id}", controller.DeleteMovie).Methods("DELETE", "OPTIONS")
	myRouter.HandleFunc("/upload", util.UploadFile)

	log.Fatal(http.ListenAndServe(":9876", myRouter))
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	log.Print("[INFO] .env file loaded")
	log.Print("[INFO] server starting...")
	handleRequests()
}
