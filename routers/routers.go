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
	router.HandleFunc("/articles", controllers.SkipCORS).Methods("OPTIONS")
	router.HandleFunc("/users", controllers.SkipCORS).Methods("OPTIONS")
	router.HandleFunc("/users/login", controllers.SkipCORS).Methods("OPTIONS")
	router.HandleFunc("/users/new", controllers.SkipCORS).Methods("OPTIONS")
	router.HandleFunc("/comments", controllers.SkipCORS).Methods("OPTIONS")
	// Router for / end point
	router.HandleFunc("/", controllers.HomePage)
	// Routers for /article... end point
	router.HandleFunc("/articles", controllers.ReturnAllArticles).Methods("GET")
	router.HandleFunc("/articles", controllers.CreateNewArticle).Methods("POST")
	router.HandleFunc("/articles/{id}", controllers.ReturnSingleArticle).Methods("GET")
	router.HandleFunc("/articles/{id}", controllers.DeleteSingleArticle).Methods("DELETE")
	router.HandleFunc("/articles/{id}", controllers.UpdateSingleArticle).Methods("PUT")
	// Routers for /user... end point
	router.HandleFunc("/users", controllers.ReturnAllUsers).Methods("GET")
	router.HandleFunc("/users/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/users/new", controllers.CreateNewUser).Methods("POST")
	router.HandleFunc("/users/{id}", controllers.ReturnSingleUser).Methods("GET")
	router.HandleFunc("/users/{id}", controllers.DeleteSingleUser).Methods("DELETE")
	router.HandleFunc("/users/{id}", controllers.UpdateSingleUser).Methods("PUT")
	// Routers for /comment... end point
	router.HandleFunc("/comments", controllers.ReturnAllComments).Methods("GET")
	router.HandleFunc("/comments", controllers.CreateNewComment).Methods("POST")
	router.HandleFunc("/comments/{id}", controllers.ReturnSingleComment).Methods("GET")
	router.HandleFunc("/comments/{id}", controllers.DeleteSingleComment).Methods("DELETE")
	router.HandleFunc("/comments/{id}", controllers.UpdateSingleComment).Methods("PUT")
	// Routers for /comment/{id}/replies... end point
	router.HandleFunc("/comments/{id}/replies", controllers.ReturnAllCommentReplies).Methods("GET")
	router.HandleFunc("/comments/{id}/replies", controllers.CreateNewCommentReply).Methods("POST")
	router.HandleFunc("/comments/{id}/replies/{rd}", controllers.ReturnSingleCommentReply).Methods("GET")
	router.HandleFunc("/comments/{id}/replies/{rd}", controllers.DeleteSingleCommentReply).Methods("DELETE")
	// Routers for /like... end point
	router.HandleFunc("/likes", controllers.ReturnAllLikes).Methods("GET")
	router.HandleFunc("/likes", controllers.CreateNewLike).Methods("POST")
	router.HandleFunc("/likes/{id}", controllers.ReturnSingleLike).Methods("GET")
	router.HandleFunc("/likes/{id}", controllers.DeleteSingleLike).Methods("DELETE")
	router.HandleFunc("/likes/{id}", controllers.UpdateSingleLike).Methods("PUT")
	// Routers for /agree... end point
	router.HandleFunc("/agrees", controllers.ReturnAllAgrees).Methods("GET")
	// Get port from environments
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	fmt.Println("Listing at: 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
