package controllers

import (
	"fmt"
	"net/http"
)

// SkipCORS : Allow CORS on all requests
func SkipCORS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: skipOPTION")
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Write([]byte("Ok!"))
}
