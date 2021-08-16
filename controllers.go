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
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gitlab.com/greenly/go-rest-api/database"
	"gitlab.com/greenly/go-rest-api/models"
)

var (
	Users    []models.User
	Articles []models.Article
	Comments []models.Comment
	db       = new(database.Database).GetDatabase()
	AppCache = new(database.Database).GetCache()
	fillFunc = func() bool {
		db.Preload("Articles").Preload("Comments").Find(&Users)
		db.Preload("Comments").Find(&Articles)
		db.Preload("Replies").Find(&Comments)
		AppCache.Set("users", Users, 24*time.Hour)
		AppCache.Set("articles", Articles, 24*time.Hour)
		AppCache.Set("comments", Comments, 24*time.Hour)
		return true
	}
	fill = fillFunc()
)

func findModel(id string, modelType string) interface{} {
	list, exist := AppCache.Get(modelType)
	if exist {
		if modelType == "users" {
			return findUser(id, list.([]models.User))
		} else if modelType == "articles" {
			return findArticle(id, list.([]models.Article))
		} else if modelType == "comments" {
			return findComment(id, list.([]models.Comment))
		} else {
			return nil
		}
	}
	return nil
}

func findUser(id string, users []models.User) interface{} {
	for _, user := range users {
		if fmt.Sprint(user.ID) == id {
			return user
		}
	}
	return nil
}

func findArticle(id string, articles []models.Article) interface{} {
	for _, article := range articles {
		if fmt.Sprint(article.ID) == id {
			return article
		}
	}
	return nil
}

func findComment(id string, comments []models.Comment) interface{} {
	for _, comment := range comments {
		if fmt.Sprint(comment.ID) == id {
			return comment
		}
	}
	return nil
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
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
	} else {
		result := fmt.Sprintf("No article found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article models.Article
	json.Unmarshal(reqBody, &article)
	fmt.Printf("Endpoint Hit: CreateNewArticle by id='%v'\n", article.ID)
	found := findModel(fmt.Sprint(article.ID), "articles")
	if found == nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		userId := r.Context().Value("user").(uint)
		article.UserID = uint(userId)
		db.Create(&article)
		db.Preload("Comments").Find(&Articles)
		db.Preload("Articles").Preload("Comments").Find(&Users)
		AppCache.Set("users", Users, 24*time.Hour)
		AppCache.Set("articles", Articles, 24*time.Hour)
		json.NewEncoder(w).Encode(article)
	} else {
		result := fmt.Sprintf("One article found by id: '%v'!", article.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func deleteSingleArticle(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findModel(fmt.Sprint(senderId), "users")
	sender := senderFound.(models.User)
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleArticle by id='%v'\n", id)
	found := findModel(id, "articles")
	if found != nil {
		article := found.(models.Article)
		if uint(senderId) != article.UserID && sender.Admin == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		db.Delete(&article)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		db.Find(&Comments)
		db.Preload("Comments").Find(&Articles)
		db.Preload("Articles").Preload("Comments").Find(&Users)
		AppCache.Set("comments", Comments, 24*time.Hour)
		AppCache.Set("articles", Articles, 24*time.Hour)
		AppCache.Set("users", Users, 24*time.Hour)
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
	sender := senderFound.(models.User)
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleArticle by id='%v'\n", id)
	found := findModel(id, "articles")
	if found != nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var article models.Article
		var reqMap map[string]string
		db.First(&article, id)
		json.Unmarshal(reqBody, &reqMap)
		if uint(senderId) != article.UserID && sender.Admin == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		article.Title = reqMap["Title"]
		article.Desc = reqMap["Descriptions"]
		article.Content = reqMap["Content"]
		db.Save(&article)
		db.Preload("Comments").Find(&Articles)
		db.Preload("Articles").Preload("Comments").Find(&Users)
		AppCache.Set("users", Users, 24*time.Hour)
		AppCache.Set("articles", Articles, 24*time.Hour)
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
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
	} else {
		result := fmt.Sprintf("No user found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user models.User
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
			sender := senderFound.(models.User)
			if sender.Admin == false {
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
		db.Preload("Articles").Preload("Comments").Find(&Users)
		AppCache.Set("users", Users, 24*time.Hour)
		// Show every things about new user exp: Hashed password, jwt token
		var tmpUser struct {
			ID        uint
			FirstName string           `json:"first_name"`
			LastName  string           `json:"last_name"`
			Email     string           `json:"email"`
			Age       string           `json:"age"`
			Username  string           `json:"username"`
			Password  string           `json:"password"`
			Admin     bool             `json:"admin"`
			Token     string           `json:"token"`
			Articles  []models.Article `json:"articles"`
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
	senderId := r.Context().Value("user").(uint)
	senderFound := findModel(fmt.Sprint(senderId), "users")
	sender := senderFound.(models.User)
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleUser by id='%v'\n", id)
	// Only owners or admins can do this
	if strconv.Itoa(int(senderId)) != id && sender.Admin == false {
		http.Error(w, "Permission Dinied!", http.StatusForbidden)
		return
	}
	found := findModel(id, "users")
	if found != nil {
		user := found.(models.User)
		db.Delete(&user)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		db.Find(&Comments)
		db.Preload("Comments").Find(&Articles)
		db.Preload("Articles").Preload("Comments").Find(&Users)
		AppCache.Set("comments", Comments, 24*time.Hour)
		AppCache.Set("articles", Articles, 24*time.Hour)
		AppCache.Set("users", Users, 24*time.Hour)
		result := user
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No user found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func updateSingleUser(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findModel(fmt.Sprint(senderId), "users")
	sender := senderFound.(models.User)
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleUser by id='%v'\n", id)
	// Only owners or admins can do this
	if strconv.Itoa(int(senderId)) != id && sender.Admin == false {
		http.Error(w, "Permission Dinied!", http.StatusForbidden)
		return
	}
	found := findModel(id, "users")
	if found != nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var user models.User
		var reqMap map[string]interface{}
		db.First(&user, id)
		json.Unmarshal(reqBody, &reqMap)
		user.FirstName = reqMap["first_name"].(string)
		user.LastName = reqMap["last_name"].(string)
		user.Email = reqMap["email"].(string)
		user.Age = reqMap["age"].(string)
		user.Username = reqMap["username"].(string)
		user.Password = reqMap["password"].(string)
		if sender.Admin == true {
			user.Admin = reqMap["admin"].(bool)
		}
		res, ok := user.Update()
		if !ok {
			http.Error(w, res, http.StatusBadRequest)
			return
		}
		db.Preload("Articles").Preload("Comments").Find(&Users)
		AppCache.Set("users", Users, 24*time.Hour)
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
			tk := &models.Token{UserId: user.ID}
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
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
	} else {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		w.WriteHeader(404)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func createNewComment(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var comment models.Comment
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
		db.Find(&Comments)
		db.Preload("Comments").Find(&Articles)
		db.Preload("Articles").Preload("Comments").Find(&Users)
		AppCache.Set("comments", Comments, 24*time.Hour)
		AppCache.Set("articles", Articles, 24*time.Hour)
		AppCache.Set("users", Users, 24*time.Hour)
		json.NewEncoder(w).Encode(comment)
	} else {
		result := fmt.Sprintf("One comment found by id: '%v'!", comment.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

func deleteSingleComment(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findModel(fmt.Sprint(senderId), "users")
	sender := senderFound.(models.User)
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleComment by id='%v'\n", id)
	found := findModel(id, "comments")
	if found != nil {
		comment := found.(models.Comment)
		if uint(senderId) != comment.UserID && sender.Admin == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		db.Delete(&comment)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		db.Find(&Comments)
		db.Preload("Comments").Find(&Articles)
		db.Preload("Articles").Preload("Comments").Find(&Users)
		AppCache.Set("comments", Comments, 24*time.Hour)
		AppCache.Set("articles", Articles, 24*time.Hour)
		AppCache.Set("users", Users, 24*time.Hour)
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
	sender := senderFound.(models.User)
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleComment by id='%v'\n", id)
	found := findModel(id, "comments")
	if found != nil {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var comment models.Comment
		var reqMap map[string]string
		db.First(&comment, id)
		json.Unmarshal(reqBody, &reqMap)
		if senderId != comment.UserID && sender.Admin == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		comment.Message = reqMap["Message"]
		db.Save(&comment)
		db.Find(&Comments)
		db.Preload("Comments").Find(&Articles)
		db.Preload("Articles").Preload("Comments").Find(&Users)
		AppCache.Set("comments", Comments, 24*time.Hour)
		AppCache.Set("articles", Articles, 24*time.Hour)
		AppCache.Set("users", Users, 24*time.Hour)
		result := comment
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)

	}
}
