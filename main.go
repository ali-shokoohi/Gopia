package main

import (
	"encoding/json"
	"fmt"
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

var Articles []Article

func findArticle(id string) []Article {
	var found []Article
	for _, article := range Articles {
		if article.Id == id {
			found = append(found, article)
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

func returnSingleArticles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeArticles by id='%v'\n", id)
	found := findArticle(id)
	if found != nil {
		result := found[0]
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

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/article", returnAllArticles)
	router.HandleFunc("/article/{id}", returnSingleArticles)
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
