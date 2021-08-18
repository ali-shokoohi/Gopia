package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var ()

// Find a object of specify model
func findObject(id string, models ...interface{}) interface{} {
	// We get only one model here, So:
	model := models[0]
	// Filter our model's objects with specify ID if It's exists
	filtered := filter("ID", id, model)
	if filtered != nil {
		return filtered[0] // Cause ID is a primarykey in table, We have a maximum of one record
	}
	return nil
}

// Filter objects in a slice (Array, List)
func filter(key string, value string, slices ...interface{}) []map[string]interface{} {
	StringList, err := json.Marshal(slices[0]) // [0] is because for we have only one ...interface{}
	if err != nil {
		panic(err)
	}
	// Convert []byte to slice of map[string]interface{}
	var list []map[string]interface{}
	err = json.Unmarshal(StringList, &list)
	if err != nil {
		panic(err)
	}
	// Search for our value
	var found []map[string]interface{}
	for _, element := range list {
		if fmt.Sprint(element[key]) == value {
			found = append(found, element)
		}
	}
	if len(found) == 0 {
		return nil
	}
	return found
}

// HomePage controller
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

// SkipCORS controller
func SkipCORS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: skipOPTION")
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Write([]byte("Ok!"))
}
