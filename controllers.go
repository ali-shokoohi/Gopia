package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type FoundModel struct {
	Index       int
	ModelObject interface{}
}

func findModel(id string, modelType string) []FoundModel {
	var found []FoundModel
	for index, model := range objectsJsonMap[modelType] {
		if fmt.Sprint(model.(map[string]interface{})["ID"]) == id {
			FoundModel := FoundModel{Index: index, ModelObject: model}
			found = append(found, FoundModel)
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
	fmt.Printf("The fucking id is: %v\n", article.ID)
	fmt.Printf("Endpoint Hit: CreateNewArticle by id='%v'\n", article.ID)
	found := findModel(fmt.Sprint(article.ID), "articles")
	if found == nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		userUrl := r.URL.User.String()
		userId, _ := strconv.ParseUint(userUrl, 10, 64)
		article.UserID = uint(userId)
		db.Create(&article)
		// Reload Users list
		db.Preload("Articles").Find(&Users)
		Articles = append(Articles, article)
		objects["articles"] = Articles
		objects["users"] = Users
		reloadObjects()
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
		// Reload Users list
		db.Preload("Articles").Find(&Users)
		objects["articles"] = Articles
		objects["users"] = Users
		reloadObjects()
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
		// Reload Users list
		db.Preload("Articles").Find(&Users)
		Articles = append(Articles, article)
		objects["articles"] = Articles
		objects["users"] = Users
		reloadObjects()
		result := article
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No article found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)

	}
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllUsers")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Users)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeUser by id='%v'\n", id)
	found := findModel(id, "users")
	if found != nil {
		result := found[0].ModelObject
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No user found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)
	fmt.Printf("Endpoint Hit: CreateNewUser by id='%v'\n", user.ID)
	found := findModel(fmt.Sprint(user.ID), "users")
	if found == nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		res, ok := user.Create()
		if !ok {
			http.Error(w, res, http.StatusBadRequest)
			return
		}
		Users = append(Users, user)
		objects["users"] = Users
		reloadObjects()
		// Show every things about new user exp: Hashed password, jwt token
		var tmpUser struct {
			ID        uint
			FirstName string    `json:"first_name"`
			LastName  string    `json:"last_name"`
			Email     string    `json:"email"`
			Age       string    `json:"age"`
			Username  string    `json:"username"`
			Password  string    `json:"password"`
			Token     string    `json:"token"`
			Articles  []Article `json:"articles"`
		}
		tmpUser.ID = user.ID
		tmpUser.FirstName = user.FirstName
		tmpUser.LastName = user.LastName
		tmpUser.Email = user.Email
		tmpUser.Username = user.Username
		tmpUser.Password = user.Password
		tmpUser.Age = user.Age
		tmpUser.Token = user.Token
		tmpUser.Articles = user.Articles
		json.NewEncoder(w).Encode(tmpUser)
	} else {
		result := fmt.Sprintf("One user found by id: '%v'!", user.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func deleteSingleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleUser by id='%v'\n", id)
	found := findModel(id, "users")
	if found != nil {
		user := found[0].ModelObject
		index := found[0].Index
		db.Delete(&user)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		Users = append(Users[:index], Users[index+1:]...)
		// Reload Articles list
		db.Find(&Articles)
		objects["users"] = Users
		objects["articles"] = Articles
		reloadObjects()
		result := user
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No user found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func updateSingleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleUser by id='%v'\n", id)
	found := findModel(id, "users")
	if found != nil {
		index := found[0].Index
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		Users = append(Users[:index], Users[index+1:]...)
		reqBody, _ := ioutil.ReadAll(r.Body)
		var user User
		var reqMap map[string]string
		db.First(&user, id)
		json.Unmarshal(reqBody, &reqMap)
		user.FirstName = reqMap["first_name"]
		user.LastName = reqMap["last_name"]
		user.Email = reqMap["email"]
		user.Age = reqMap["age"]
		user.Username = reqMap["username"]
		user.Password = reqMap["password"]
		res, ok := user.Update()
		if !ok {
			http.Error(w, res, http.StatusBadRequest)
			return
		}
		Users = append(Users, user)
		objects["users"] = Users
		reloadObjects()
		result := user
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No user found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)

	}
}
