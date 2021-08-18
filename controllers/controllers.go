package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/greenly/go-rest-api/models"
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

// ReturnAllComments controller
func ReturnAllComments(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.RawQuery
	// If rawQuery is exists in request. ex: UserID=1
	if len(rawQuery) > 0 {
		fmt.Printf("Endpoint Hit: ReturnAllComments by '%s'\n", rawQuery)
		cut := strings.Split(rawQuery, "=")
		key, value := cut[0], cut[1]
		found := filter(key, value, models.Comments)
		if found == nil {
			result := fmt.Sprintf("No comment found by '%s': '%s'!", key, value)
			w.WriteHeader(404)
			http.Error(w, result, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
		return
	}
	fmt.Println("Endpoint Hit: returnAllComments")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Comments)
}

// ReturnSingleComment controller
func ReturnSingleComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeComment by id='%v'\n", id)
	found := findObject(id, models.Comments)
	if found != nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
	} else {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// CreateNewComment controller
func CreateNewComment(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var comment models.Comment
	json.Unmarshal(reqBody, &comment)
	fmt.Printf("Endpoint Hit: CreateNewComment by id='%v'\n", comment.ID)
	found := findObject(fmt.Sprint(comment.ID), models.Comments)
	if found == nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		userId := r.Context().Value("user").(uint)
		comment.UserID = uint(userId)
		models.DB.Create(&comment)
		// Reload Users list
		models.DB.Find(&models.Comments)
		models.DB.Preload("Comments").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		json.NewEncoder(w).Encode(comment)
	} else {
		result := fmt.Sprintf("One comment found by id: '%v'!", comment.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// DeleteSingleComment controller
func DeleteSingleComment(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findObject(fmt.Sprint(senderId), models.Users)
	sender := senderFound.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleComment by id='%v'\n", id)
	found := findObject(id, models.Comments)
	if found != nil {
		comment := found.(map[string]interface{})
		if uint(senderId) != comment["UserID"] && sender["Admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		models.DB.Delete(&models.Comment{}, comment["ID"])
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		models.DB.Find(&models.Comments)
		models.DB.Preload("Comments").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		result := comment
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// UpdateSingleComment controller
func UpdateSingleComment(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findObject(fmt.Sprint(senderId), models.Users)
	sender := senderFound.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleComment by id='%v'\n", id)
	found := findObject(id, models.Comments)
	if found != nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var comment models.Comment
		var reqMap map[string]string
		models.DB.First(&comment, id)
		json.Unmarshal(reqBody, &reqMap)
		if senderId != comment.UserID && sender["Admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		comment.Message = reqMap["Message"]
		models.DB.Save(&comment)
		models.DB.Find(&models.Comments)
		models.DB.Preload("Comments").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		result := comment
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)

	}
}
