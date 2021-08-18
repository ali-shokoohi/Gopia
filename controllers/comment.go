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

// ReturnAllComments : Return all comments with or without raw query in request
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

// ReturnSingleComment : Return a comment article by ID
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

// CreateNewComment : Create a new comment object
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
		models.DB.Preload("Replies").Find(&models.Comments)
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

// DeleteSingleComment : Delete a single comment object from database by ID
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
		models.DB.Preload("Replies").Find(&models.Comments)
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

// UpdateSingleComment : Update and change a single comment in database via ID
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
		models.DB.Preload("Replies").Find(&models.Comments)
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
