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

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"Descriptions"`
	Content string `json:"Content"`
}

type FoundArticle struct {
	Index         int     `json:"Index"`
	ArticleObject Article `json:"Article"`
}

var Articles []Article

func findArticle(id string) []FoundArticle {
	var found []FoundArticle
	for index, article := range Articles {
		if article.Id == id {
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
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeArticle by id='%v'\n", id)
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
	fmt.Printf("Endpoint Hit: CreateNewArticle by id='%v'\n", article.Id)
	found := findArticle(string(article.Id))
	if found == nil {
		w.WriteHeader(200)
		Articles = append(Articles, article)
		json.NewEncoder(w).Encode(article)
	} else {
		result := HttpError{
			StatusCode: 400,
			Message:    fmt.Sprintf("One article found by id: '%v'!", article.Id),
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(result)
	}
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteArticle by id='%v'\n", id)
	found := findArticle(id)
	if found != nil {
		article := found[0].ArticleObject
		index := found[0].Index
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

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/article", returnAllArticles).Methods("GET")
	router.HandleFunc("/article", createNewArticle).Methods("POST")
	router.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET")
	router.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Articles = []Article{
		Article{Id: "0", Title: "1984", Desc: "Article of 1984 book", Content: "This book is wonderful"},
		Article{Id: "1", Title: "Homo sapiens", Desc: "Article of Homo sapiens book", Content: "This book is so useful"},
	}
	handleRequests()
}
