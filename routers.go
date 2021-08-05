package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	// Enable CORS for all endpoints
	router.Use(CORSMiddleWare)
	// Show request URL
	router.Use(urlMiddleWare)
	// Basic Authentication middleware
	//router.Use(authMiddleWare)
	// JWT Authentication middleware
	router.Use(jwtMiddleWare)
	// Route for all end points to skip OPTIONS method for CORS
	router.HandleFunc("/article", skipCORS).Methods("OPTIONS")
	router.HandleFunc("/user", skipCORS).Methods("OPTIONS")
	router.HandleFunc("/user/login", skipCORS).Methods("OPTIONS")
	router.HandleFunc("/user/new", skipCORS).Methods("OPTIONS")
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
	router.HandleFunc("/user/login", loginUser).Methods("POST")
	router.HandleFunc("/user/new", createNewUser).Methods("POST")
	router.HandleFunc("/user/{id}", returnSingleUser).Methods("GET")
	router.HandleFunc("/user/{id}", deleteSingleUser).Methods("DELETE")
	router.HandleFunc("/user/{id}", updateSingleUser).Methods("PUT")
	// Routers for /comment... end point
	router.HandleFunc("/comment", returnAllComments).Methods("GET")
	router.HandleFunc("/comment", createNewComment).Methods("POST")
	router.HandleFunc("/comment/{id}", returnSingleComment).Methods("GET")
	router.HandleFunc("/comment/{id}", deleteSingleComment).Methods("DELETE")
	router.HandleFunc("/comment/{id}", updateSingleComment).Methods("PUT")
	// Get port from environments
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	fmt.Println("Listing at: 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
