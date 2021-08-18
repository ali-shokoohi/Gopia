package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"gitlab.com/greenly/go-rest-api/models"
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
