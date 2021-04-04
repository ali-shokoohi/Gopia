package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type FoundModel struct {
	Index       int
	ModelObject interface{}
}

func findModel(id string, modelType string) []FoundModel {
	var found []FoundModel
	// Will be better later!
	// TODO: Make this better to detect types automatic...
	// Like: objects[modelType].([]interface{}) with fields for all types
	switch modelType {
	case "articles":
		for index, model := range objects[modelType].([]Article) {
			if fmt.Sprint(model.ID) == id {
				FoundModel := FoundModel{Index: index, ModelObject: model}
				found = append(found, FoundModel)
				break
			}
		}
	}
	return found
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeArticle by id='%v'\n", id)
	found := findModel(id, "articles")
	if found != nil {
		result := found[0].ModelObject
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No article found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	fmt.Printf("Endpoint Hit: CreateNewArticle by id='%v'\n", article.ID)
	found := findModel(fmt.Sprint(article.ID), "articles")
	if found == nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		db.Create(&article)
		Articles = append(Articles, article)
		json.NewEncoder(w).Encode(article)
	} else {
		result := fmt.Sprintf("One article found by id: '%v'!", article.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func deleteSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleArticle by id='%v'\n", id)
	found := findModel(id, "articles")
	if found != nil {
		article := found[0].ModelObject
		index := found[0].Index
		db.Delete(&article)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		Articles = append(Articles[:index], Articles[index+1:]...)
		result := article
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No article found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func updateSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleArticle by id='%v'\n", id)
	found := findModel(id, "articles")
	if found != nil {
		index := found[0].Index
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		Articles = append(Articles[:index], Articles[index+1:]...)
		reqBody, _ := ioutil.ReadAll(r.Body)
		var article Article
		var reqMap map[string]string
		db.First(&article, id)
		json.Unmarshal(reqBody, &reqMap)
		article.Title = reqMap["Title"]
		article.Desc = reqMap["Descriptions"]
		article.Content = reqMap["Content"]
		db.Save(&article)
		Articles = append(Articles, article)
		result := article
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No article found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)

	}
}
