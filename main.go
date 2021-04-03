package main

import (
	"fmt"
)

var Articles []Article

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	models = perpareModels()
	db.Find(&Articles)
	handleRequests()
}
