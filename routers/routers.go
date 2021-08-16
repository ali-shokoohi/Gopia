package routers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gitlab.com/greenly/go-rest-api/controllers"
	"gitlab.com/greenly/go-rest-api/utils/middlewares"
)

// HandleRequests for handle all requests
func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	// Enable CORS for all endpoints
	router.Use(middlewares.CORSMiddleWare)
	// Show request URL
	router.Use(middlewares.URLMiddleWare)
	// Basic Authentication middleware
	//router.Use(authMiddleWare)
	// JWT Authentication middleware
	router.Use(middlewares.JWTMiddleWare)
	// Route for all end points to skip OPTIONS method for CORS
	router.HandleFunc("/article", controllers.SkipCORS).Methods("OPTIONS")
	router.HandleFunc("/user", controllers.SkipCORS).Methods("OPTIONS")
	router.HandleFunc("/user/login", controllers.SkipCORS).Methods("OPTIONS")
	router.HandleFunc("/user/new", controllers.SkipCORS).Methods("OPTIONS")
	// Router for / end point
	router.HandleFunc("/", controllers.HomePage)
	// Routers for /article... end point
	router.HandleFunc("/article", controllers.ReturnAllArticles).Methods("GET")
	router.HandleFunc("/article", controllers.CreateNewArticle).Methods("POST")
	router.HandleFunc("/article/{id}", controllers.ReturnSingleArticle).Methods("GET")
	router.HandleFunc("/article/{id}", controllers.DeleteSingleArticle).Methods("DELETE")
	router.HandleFunc("/article/{id}", controllers.UpdateSingleArticle).Methods("PUT")
	// Routers for /user... end point
	router.HandleFunc("/user", controllers.ReturnAllUsers).Methods("GET")
	router.HandleFunc("/user/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/user/new", controllers.CreateNewUser).Methods("POST")
	router.HandleFunc("/user/{id}", controllers.ReturnSingleUser).Methods("GET")
	router.HandleFunc("/user/{id}", controllers.DeleteSingleUser).Methods("DELETE")
	router.HandleFunc("/user/{id}", controllers.UpdateSingleUser).Methods("PUT")
	// Routers for /comment... end point
	router.HandleFunc("/comment", controllers.ReturnAllComments).Methods("GET")
	router.HandleFunc("/comment", controllers.CreateNewComment).Methods("POST")
	router.HandleFunc("/comment/{id}", controllers.ReturnSingleComment).Methods("GET")
	router.HandleFunc("/comment/{id}", controllers.DeleteSingleComment).Methods("DELETE")
	router.HandleFunc("/comment/{id}", controllers.UpdateSingleComment).Methods("PUT")
	// Get port from environments
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	fmt.Println("Listing at: 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
