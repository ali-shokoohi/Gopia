package main

import (
	"fmt"
)

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	models, objects = perpareModels()
	handleRequests()
}
