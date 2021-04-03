package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(urlMiddleWare)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/article", returnAllArticles).Methods("GET")
	router.HandleFunc("/article", createNewArticle).Methods("POST")
	router.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET")
	router.HandleFunc("/article/{id}", deleteSingleArticle).Methods("DELETE")
	router.HandleFunc("/article/{id}", updateSingleArticle).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8090", router))
}
