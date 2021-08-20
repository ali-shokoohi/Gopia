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

// ReturnAllCommentReplies : Return all of the comment's replies with or without raw query in request
func ReturnAllCommentReplies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	found := findObject(id, models.Comments)
	if found == nil {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
		return
	}
	replies := found.(map[string]interface{})["replies"]
	// If rawQuery is exists in request. ex: UserID=1
	rawQuery := r.URL.RawQuery
	if len(rawQuery) > 0 {
		fmt.Printf("Endpoint Hit: returnAllCommentReplies of comment: %v by '%s'\n", id, rawQuery)
		cut := strings.Split(rawQuery, "=")
		key, value := cut[0], cut[1]
		found := filter(key, value, replies)
		if found == nil {
			result := fmt.Sprintf("No replies found by '%s': '%s'!", key, value)
			w.WriteHeader(404)
			http.Error(w, result, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
		return
	}
	fmt.Printf("Endpoint Hit: returnAllCommentReplies of comment: %v\n", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}

// ReturnSingleCommentReply : Return a comment's reply by ID
func ReturnSingleCommentReply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeCommentReply by id='%v'\n", id)
	found := findObject(id, models.Comments)
	if found == nil {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
		return
	}
	replies := found.(map[string]interface{})["replies"]
	// Search for a reply via reply ID
	rd := vars["rd"]
	filterred := filter("ID", rd, replies)
	if filterred == nil {
		result := fmt.Sprintf("No reply found by id: %v!", rd)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
		return
	}
	reply := filterred[0] // [0] because ID field is a primarykey in database
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

// CreateNewCommentReply : Create a new reply of a comment object
func CreateNewCommentReply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeCommentReply by id='%v'\n", id)
	found := findObject(id, models.Comments)
	if found == nil {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
		return
	}
	var comment models.Comment
	models.DB.First(&comment, found.(map[string]interface{})["ID"])
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reply models.Comment
	json.Unmarshal(reqBody, &reply)
	fmt.Printf("Endpoint Hit: CreateNewCommentReply by id='%v'\n", reply.ID)
	found = findObject(fmt.Sprint(reply.ID), models.Comments)
	if found == nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		userId := r.Context().Value("user").(uint)
		reply.UserID = uint(userId)
		// Create reply
		models.DB.Create(&reply)
		// Append reply to its comment
		comment.Replies = append(comment.Replies, &reply)
		models.DB.Save(&comment)
		// Reload Users list
		models.DB.Preload(clause.Associations).Find(&models.Comments)
		models.DB.Preload(clause.Associations).Find(&models.Articles)
		models.DB.Preload(clause.Associations).Find(&models.Users)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		json.NewEncoder(w).Encode(comment)
	} else {
		result := fmt.Sprintf("One comment found by id: '%v'!", comment.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// DeleteSingleCommentReply : Delete a single comment's reply object from database by ID
func DeleteSingleCommentReply(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findObject(fmt.Sprint(senderId), models.Users)
	sender := senderFound.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	rd := vars["rd"]
	fmt.Printf("Endpoint Hit: deleteSingleCommentReply by id='%v'\n", rd)
	found := findObject(id, models.Comments)
	if found != nil {
		comment := found.(map[string]interface{})
		if uint(senderId) != comment["UserID"] && sender["Admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		replies := comment["replies"]
		filterred := filter("ID", rd, replies)
		if filterred == nil {
			result := fmt.Sprintf("No reply found by id '%v' in comment's replies '%v'!", rd, id)
			http.Error(w, result, http.StatusBadRequest)
			return
		}
		reply := filterred[0] // [0] because ID field is a primarykey in database
		models.DB.Delete(&models.Comment{}, reply["ID"])
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		models.DB.Preload(clause.Associations).Find(&models.Comments)
		models.DB.Preload(clause.Associations).Find(&models.Articles)
		models.DB.Preload(clause.Associations).Find(&models.Users)
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
