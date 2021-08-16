package main

import (
	"fmt"

	"gitlab.com/greenly/go-rest-api/models"
	"gitlab.com/greenly/go-rest-api/routers"
)

func main() {
	fmt.Println("Go-Rest-API")
	new(models.Model).PerpareModels()
	routers.HandleRequests()
}
