package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HttpError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// type Summarizer interface {
// 	summarize() string
// }

// func (a *Article) summarize() string {
// 	return fmt.Sprintf("%s: %s", a.Title, a.Desc)
// }

type FoundArticle struct {
	Index         int     `json:"Index"`
	ArticleObject Article `json:"Article"`
}

var Articles []Article

func findArticle(id string) []FoundArticle {
	var found []FoundArticle
	for index, article := range Articles {
		if fmt.Sprint(article.ID) == id {
			foundArticle := FoundArticle{Index: index, ArticleObject: article}
			found = append(found, foundArticle)
			break
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
	w.Header().Set("Content-Type", "application/json")
	found := findArticle(id)
	if found != nil {
		result := found[0].ArticleObject
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
	} else {
		result := HttpError{
			StatusCode: 404,
			Message:    fmt.Sprintf("No article found by id: '%v'!", id),
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(result)
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	fmt.Printf("Endpoint Hit: CreateNewArticle by id='%v'\n", article.ID)
	w.Header().Set("Content-Type", "application/json")
	found := findArticle(fmt.Sprint(article.ID))
	if found == nil {
		w.WriteHeader(200)
		db.Create(&article)
		Articles = append(Articles, article)
		json.NewEncoder(w).Encode(article)
	} else {
		result := HttpError{
			StatusCode: 400,
			Message:    fmt.Sprintf("One article found by id: '%v'!", article.ID),
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(result)
	}
}

func deleteSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleArticle by id='%v'\n", id)
	w.Header().Set("Content-Type", "application/json")
	found := findArticle(id)
	if found != nil {
		article := found[0].ArticleObject
		index := found[0].Index
		db.Delete(&article)
		w.WriteHeader(200)
		Articles = append(Articles[:index], Articles[index+1:]...)
		result := article
		json.NewEncoder(w).Encode(result)
	} else {
		result := HttpError{
			StatusCode: 404,
			Message:    fmt.Sprintf("No article found by id: '%v'!", id),
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(result)
	}
}

func updateSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleArticle by id='%v'\n", id)
	w.Header().Set("Content-Type", "application/json")
	found := findArticle(id)
	if found != nil {
		index := found[0].Index
		w.WriteHeader(200)
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
		result := HttpError{
			StatusCode: 404,
			Message:    fmt.Sprintf("No article found by id: '%v'!", id),
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(result)
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/article", returnAllArticles).Methods("GET")
	router.HandleFunc("/article", createNewArticle).Methods("POST")
	router.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET")
	router.HandleFunc("/article/{id}", deleteSingleArticle).Methods("DELETE")
	router.HandleFunc("/article/{id}", updateSingleArticle).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	models = perpareModels()
	fmt.Println("Rest API v2.0 - Mux Routers")
	db.Find(&Articles)
	handleRequests()
}
