package controllers

import (
	"fmt"
	"net/http"
)

// HomePage : Say welcome in root end point '/'
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
