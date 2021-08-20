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
	"gorm.io/gorm/clause"
)

// ReturnAllAgrees - Return all agrees with or without raw query in request
func ReturnAllAgrees(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.RawQuery
	// If rawQuery is exists in request. ex: UserID=1
	if len(rawQuery) > 0 {
		fmt.Printf("Endpoint Hit: ReturnAllAgrees by '%s'\n", rawQuery)
		cut := strings.Split(rawQuery, "=")
		key, value := cut[0], cut[1]
		found := filter(key, value, models.Agrees)
		if found == nil {
			result := fmt.Sprintf("No agree found by '%s': '%s'!", key, value)
			w.WriteHeader(404)
			http.Error(w, result, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
		return
	}
	fmt.Println("Endpoint Hit: returnAllAgrees")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Agrees)
}

// ReturnSingleAgree - Return a agree article by ID
func ReturnSingleAgree(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeAgree by id='%v'\n", id)
	found := findObject(id, models.Agrees)
	if found != nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
	} else {
		result := fmt.Sprintf("No agree found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// CreateNewAgree - Create a new agree object
func CreateNewAgree(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var agree models.Agree
	json.Unmarshal(reqBody, &agree)
	fmt.Printf("Endpoint Hit: CreateNewAgree by id='%v'\n", agree.ID)
	found := findObject(fmt.Sprint(agree.ID), models.Agrees)
	if found == nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		userId := r.Context().Value("user").(uint)
		agree.UserID = uint(userId)
		models.DB.Create(&agree)
		// Reload Users list
		models.DB.Find(&models.Agrees)
		models.DB.Preload(clause.Associations).Find(&models.Comments)
		models.DB.Preload(clause.Associations).Find(&models.Users)
		models.AppCache.Set("agrees", models.Agrees, 24*time.Hour)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		json.NewEncoder(w).Encode(agree)
	} else {
		result := fmt.Sprintf("One agree found by id: '%v'!", agree.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// DeleteSingleAgree - Delete a single agree object from database by ID
func DeleteSingleAgree(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findObject(fmt.Sprint(senderId), models.Users)
	sender := senderFound.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleAgree by id='%v'\n", id)
	found := findObject(id, models.Agrees)
	if found != nil {
		agree := found.(map[string]interface{})
		if uint(senderId) != agree["UserID"] && sender["Admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		models.DB.Delete(&models.Agree{}, agree["ID"])
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		models.DB.Find(&models.Agrees)
		models.DB.Preload(clause.Associations).Find(&models.Comments)
		models.DB.Preload(clause.Associations).Find(&models.Users)
		models.AppCache.Set("agrees", models.Agrees, 24*time.Hour)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		result := agree
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No agree found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// UpdateSingleAgree - Update and change a single agree in database via ID
func UpdateSingleAgree(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findObject(fmt.Sprint(senderId), models.Users)
	sender := senderFound.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleAgree by id='%v'\n", id)
	found := findObject(id, models.Agrees)
	if found != nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var agree models.Agree
		var reqMap map[string]interface{}
		models.DB.First(&agree, id)
		json.Unmarshal(reqBody, &reqMap)
		fmt.Println(reqMap)
		if senderId != agree.UserID && sender["Admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		agree.CommentID = uint(reqMap["CommentID"].(float64))
		models.DB.Save(&agree)
		models.DB.Find(&models.Agrees)
		models.DB.Preload(clause.Associations).Find(&models.Comments)
		models.DB.Preload(clause.Associations).Find(&models.Users)
		models.AppCache.Set("agrees", models.Agrees, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		result := agree
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No agree found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)

	}
}
