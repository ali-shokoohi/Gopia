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

// ReturnAllLikes - Return all likes with or without raw query in request
func ReturnAllLikes(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.RawQuery
	// If rawQuery is exists in request. ex: UserID=1
	if len(rawQuery) > 0 {
		fmt.Printf("Endpoint Hit: ReturnAllLikes by '%s'\n", rawQuery)
		cut := strings.Split(rawQuery, "=")
		key, value := cut[0], cut[1]
		found := filter(key, value, models.Likes)
		if found == nil {
			result := fmt.Sprintf("No like found by '%s': '%s'!", key, value)
			w.WriteHeader(404)
			http.Error(w, result, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
		return
	}
	fmt.Println("Endpoint Hit: returnAllLikes")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Likes)
}

// ReturnSingleLike - Return a like article by ID
func ReturnSingleLike(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeLike by id='%v'\n", id)
	found := findObject(id, models.Likes)
	if found != nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
	} else {
		result := fmt.Sprintf("No like found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// CreateNewLike - Create a new like object
func CreateNewLike(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var like models.Like
	json.Unmarshal(reqBody, &like)
	fmt.Printf("Endpoint Hit: CreateNewLike by id='%v'\n", like.ID)
	found := findObject(fmt.Sprint(like.ID), models.Likes)
	if found == nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		userId := r.Context().Value("user").(uint)
		like.UserID = uint(userId)
		models.DB.Create(&like)
		// Reload Users list
		models.DB.Find(&models.Likes)
		models.DB.Preload("Comments").Preload("Likes").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Preload("Likes").Find(&models.Users)
		models.AppCache.Set("likes", models.Likes, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		json.NewEncoder(w).Encode(like)
	} else {
		result := fmt.Sprintf("One like found by id: '%v'!", like.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// DeleteSingleLike - Delete a single like object from database by ID
func DeleteSingleLike(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findObject(fmt.Sprint(senderId), models.Users)
	sender := senderFound.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleLike by id='%v'\n", id)
	found := findObject(id, models.Likes)
	if found != nil {
		like := found.(map[string]interface{})
		if uint(senderId) != like["UserID"] && sender["Admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		models.DB.Delete(&models.Like{}, like["ID"])
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		models.DB.Find(&models.Likes)
		models.DB.Preload("Comments").Preload("Likes").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Preload("Likes").Find(&models.Users)
		models.AppCache.Set("likes", models.Likes, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		result := like
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No like found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)
	}
}
