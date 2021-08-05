package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
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

func skipCORS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: skipOPTION")
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Write([]byte("Ok!"))
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
		userId := r.Context().Value("user").(uint)
		article.UserID = uint(userId)
		db.Create(&article)
		// Reload Users list
		db.Preload("Articles").Preload("Comments").Find(&Users)
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
	senderId := r.Context().Value("user").(uint)
	senderFound := findModel(fmt.Sprint(senderId), "users")
	sender := senderFound[0].ModelObject.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleArticle by id='%v'\n", id)
	found := findModel(id, "articles")
	if found != nil {
		article := found[0].ModelObject.(map[string]interface{})
		index := found[0].Index
		if int(senderId) != int(article["UserID"].(float64)) && sender["admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		db.Delete(&article)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		Articles = append(Articles[:index], Articles[index+1:]...)
		// Reload Users list
		db.Preload("Articles").Preload("Comments").Find(&Users)
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
	senderId := r.Context().Value("user").(uint)
	senderFound := findModel(fmt.Sprint(senderId), "users")
	sender := senderFound[0].ModelObject.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleArticle by id='%v'\n", id)
	found := findModel(id, "articles")
	if found != nil {
		index := found[0].Index
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var article Article
		var reqMap map[string]string
		db.First(&article, id)
		json.Unmarshal(reqBody, &reqMap)
		if senderId != article.UserID && sender["admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		article.Title = reqMap["Title"]
		article.Desc = reqMap["Descriptions"]
		article.Content = reqMap["Content"]
		Articles = append(Articles[:index], Articles[index+1:]...)
		db.Save(&article)
		// Reload Users list
		db.Preload("Articles").Preload("Comments").Find(&Users)
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
		// Only admins can create super users!
		senderId := r.Context().Value("user")
		if senderId == nil {
			user.Admin = false
		} else {
			senderFound := findModel(fmt.Sprint(senderId.(uint)), "users")
			sender := senderFound[0].ModelObject.(map[string]interface{})
			if sender["admin"] == false {
				user.Admin = false
			}
		}
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
			Admin     bool      `json:"admin"`
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
		tmpUser.Admin = user.Admin
		tmpUser.Token = user.Token
		tmpUser.Articles = user.Articles
		json.NewEncoder(w).Encode(tmpUser)
	} else {
		result := fmt.Sprintf("One user found by id: '%v'!", user.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func deleteSingleUser(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(int)
	senderFound := findModel(fmt.Sprint(senderId), "users")
	sender := senderFound[0].ModelObject.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleUser by id='%v'\n", id)
	// Only owners or admins can do this
	if strconv.Itoa(senderId) != id && sender["admin"] == false {
		http.Error(w, "Permission Dinied!", http.StatusForbidden)
		return
	}
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
	senderId := r.Context().Value("user").(int)
	senderFound := findModel(fmt.Sprint(senderId), "users")
	sender := senderFound[0].ModelObject.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleUser by id='%v'\n", id)
	// Only owners or admins can do this
	if strconv.Itoa(senderId) != id && sender["admin"] == false {
		http.Error(w, "Permission Dinied!", http.StatusForbidden)
		return
	}
	found := findModel(id, "users")
	if found != nil {
		index := found[0].Index
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		Users = append(Users[:index], Users[index+1:]...)
		reqBody, _ := ioutil.ReadAll(r.Body)
		var user User
		var reqMap map[string]interface{}
		db.First(&user, id)
		json.Unmarshal(reqBody, &reqMap)
		user.FirstName = reqMap["first_name"].(string)
		user.LastName = reqMap["last_name"].(string)
		user.Email = reqMap["email"].(string)
		user.Age = reqMap["age"].(string)
		user.Username = reqMap["username"].(string)
		user.Password = reqMap["password"].(string)
		if sender["admin"] == true {
			user.Admin = reqMap["admin"].(bool)
		}
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

func loginUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reqMap map[string]string
	json.Unmarshal(reqBody, &reqMap)
	reqUser := reqMap["username"]
	reqPass := reqMap["password"]
	hasher := md5.New()
	hasher.Write([]byte(reqPass))
	hashPass := hex.EncodeToString(hasher.Sum(nil))
	for _, user := range Users {
		if user.Username == reqUser && user.Password == hashPass {
			// Set user in request
			tk := &Token{UserId: user.ID}
			token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
			tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
			user.Token = tokenString //Store the token in the response
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			var result struct {
				Status string `json:"status"`
				Token  string `json:"token"`
			}
			result.Status = "success"
			result.Token = user.Token
			json.NewEncoder(w).Encode(result)
			return
		}
	}
	http.Error(w, "Access Dinied!", http.StatusForbidden)
}

func returnAllComments(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllComments")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Comments)
}

func returnSingleComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeComment by id='%v'\n", id)
	found := findModel(id, "comments")
	if found != nil {
		result := found[0].ModelObject
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func createNewComment(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var comment Comment
	json.Unmarshal(reqBody, &comment)
	fmt.Printf("Endpoint Hit: CreateNewComment by id='%v'\n", comment.ID)
	found := findModel(fmt.Sprint(comment.ID), "comments")
	if found == nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		userId := r.Context().Value("user").(uint)
		comment.UserID = uint(userId)
		db.Create(&comment)
		// Reload Users list
		db.Preload("Articles").Preload("Comments").Find(&Users)
		db.Preload("Comments").Find(&Articles)
		Comments = append(Comments, comment)
		objects["comments"] = Comments
		objects["articles"] = Articles
		objects["users"] = Users
		reloadObjects()
		json.NewEncoder(w).Encode(comment)
	} else {
		result := fmt.Sprintf("One comment found by id: '%v'!", comment.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func deleteSingleComment(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findModel(fmt.Sprint(senderId), "users")
	sender := senderFound[0].ModelObject.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleComment by id='%v'\n", id)
	found := findModel(id, "comments")
	if found != nil {
		comment := found[0].ModelObject.(map[string]interface{})
		index := found[0].Index
		if int(senderId) != int(comment["UserID"].(float64)) && sender["admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		db.Delete(&comment)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		Comments = append(Comments[:index], Comments[index+1:]...)
		// Reload Users list
		db.Preload("Articles").Preload("Comments").Find(&Users)
		db.Preload("Comments").Find(&Articles)
		objects["comments"] = Comments
		objects["articles"] = Articles
		objects["users"] = Users
		reloadObjects()
		result := comment
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func updateSingleComment(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findModel(fmt.Sprint(senderId), "users")
	sender := senderFound[0].ModelObject.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleComment by id='%v'\n", id)
	found := findModel(id, "comments")
	if found != nil {
		index := found[0].Index
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var comment Comment
		var reqMap map[string]string
		db.First(&comment, id)
		json.Unmarshal(reqBody, &reqMap)
		if senderId != comment.UserID && sender["admin"] == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		comment.Message = reqMap["Message"]
		Comments = append(Comments[:index], Comments[index+1:]...)
		db.Save(&comment)
		// Reload Users list
		db.Preload("Articles").Preload("Comments").Find(&Users)
		db.Preload("Comments").Find(&Articles)
		Comments = append(Comments, comment)
		objects["comments"] = Comments
		objects["articles"] = Articles
		objects["users"] = Users
		reloadObjects()
		result := comment
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)

	}
}
