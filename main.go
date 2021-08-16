package main

import (
	"fmt"

	"gitlab.com/greenly/go-rest-api/models"
)

func main() {
	fmt.Println("Go-Rest-API")
	new(models.Model).PerpareModels()
	handleRequests()
}
