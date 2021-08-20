package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/greenly/go-rest-api/models"
	"gitlab.com/greenly/go-rest-api/routers"
)

func main() {
	fmt.Println("Go-Rest-API")
	new(models.Model).PerpareModels()
	router := routers.CreateRouter()
	// Get port from HandleRequestsenvironments
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	fmt.Println("Listing at: 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
