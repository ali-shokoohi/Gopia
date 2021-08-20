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

// ReturnAllArticles : Return all of articles with or without raw query in request
func ReturnAllArticles(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.RawQuery
	// If rawQuery is exists in request. ex: UserID=1
	if len(rawQuery) > 0 {
		fmt.Printf("Endpoint Hit: ReturnAllArticles by '%s'\n", rawQuery)
		cut := strings.Split(rawQuery, "=")
		key, value := cut[0], cut[1]
		found := filter(key, value, models.Articles)
		if found == nil {
			result := fmt.Sprintf("No article found by '%s': '%s'!", key, value)
			w.WriteHeader(404)
			http.Error(w, result, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
		return
	}
	fmt.Println("Endpoint Hit: returnAllArticles")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Articles)
}

// ReturnSingleArticle : Return a single article by ID
func ReturnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeArticle by id='%v'\n", id)
	found := findObject(id, models.Articles)
	if found != nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
	} else {
		result := fmt.Sprintf("No article found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// CreateNewArticle : Create a new article object
func CreateNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article models.Article
	json.Unmarshal(reqBody, &article)
	fmt.Printf("Endpoint Hit: CreateNewArticle by id='%v'\n", article.ID)
	found := findObject(fmt.Sprint(article.ID), models.Articles)
	if found == nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		userId := r.Context().Value("user").(uint)
		article.UserID = uint(userId)
		models.DB.Create(&article)
		models.DB.Preload("Comments").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		json.NewEncoder(w).Encode(article)
	} else {
		result := fmt.Sprintf("One article found by id: '%v'!", article.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// DeleteSingleArticle : Delete a single article object from database by ID
func DeleteSingleArticle(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findObject(fmt.Sprint(senderId), models.Users)
	sender := senderFound.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleArticle by id='%v'\n", id)
	found := findObject(id, models.Articles)
	if found != nil {
		article := found.(map[string]interface{})
		if uint(senderId) != article["UserID"] && sender["Admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		models.DB.Delete(&models.Article{}, article["ID"])
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		models.DB.Preload("Replies").Find(&models.Comments)
		models.DB.Preload("Comments").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		result := article
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No article found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// UpdateSingleArticle : Update and change a single article in database via ID
func UpdateSingleArticle(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findObject(fmt.Sprint(senderId), models.Users)
	sender := senderFound.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleArticle by id='%v'\n", id)
	found := findObject(id, models.Articles)
	if found != nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var article models.Article
		var reqMap map[string]string
		models.DB.First(&article, id)
		json.Unmarshal(reqBody, &reqMap)
		if uint(senderId) != article.UserID && sender["Admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		article.Title = reqMap["Title"]
		article.Desc = reqMap["Descriptions"]
		article.Content = reqMap["Content"]
		models.DB.Save(&article)
		models.DB.Preload("Comments").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		result := article
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No article found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)

	}
}
