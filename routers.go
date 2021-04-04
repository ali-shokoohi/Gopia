package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	// Show request URL
	router.Use(urlMiddleWare)
	// Router for / end point
	router.HandleFunc("/", homePage)
	// Routers for /article... end point
	router.HandleFunc("/article", returnAllArticles).Methods("GET")
	router.HandleFunc("/article", createNewArticle).Methods("POST")
	router.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET")
	router.HandleFunc("/article/{id}", deleteSingleArticle).Methods("DELETE")
	router.HandleFunc("/article/{id}", updateSingleArticle).Methods("PUT")
	// Routers for /user... end point
	router.HandleFunc("/user", returnAllUsers).Methods("GET")
	router.HandleFunc("/user", createNewUser).Methods("POST")
	router.HandleFunc("/user/{id}", returnSingleUser).Methods("GET")
	router.HandleFunc("/user/{id}", deleteSingleUser).Methods("DELETE")
	router.HandleFunc("/user/{id}", updateSingleUser).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8090", router))
}
