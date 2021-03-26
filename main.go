package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"Descriptions"`
	Content string `json:"Content"`
}

var Articles []Article

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
	for _, article := range Articles {
		if article.Id == id {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/all", returnAllArticles)
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
