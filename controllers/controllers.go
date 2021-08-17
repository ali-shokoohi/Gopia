package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gitlab.com/greenly/go-rest-api/models"
)

var ()

func findModel(id string, modelType string) interface{} {
	list, exist := models.AppCache.Get(modelType)
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

// Filter objects in a slice (Array, List)
func filter(key string, value string, slices ...interface{}) []map[string]interface{} {
	StringList, err := json.Marshal(slices[0]) // [0] is because for we have only one ...interface{}
	if err != nil {
		panic(err)
	}
	// Convert []byte to slice of []map[string]interface{}
	var list []map[string]interface{}
	err = json.Unmarshal(StringList, &list)
	if err != nil {
		panic(err)
	}
	// Search for our value
	var found []map[string]interface{}
	for _, element := range list {
		if fmt.Sprint(element[key]) == value {
			found = append(found, element)
		}
	}
	if len(found) == 0 {
		return nil
	}
	return found
}

// HomePage controller
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

// SkipCORS controller
func SkipCORS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: skipOPTION")
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Write([]byte("Ok!"))
}

// ReturnAllArticles controller
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

// ReturnSingleArticle controller
func ReturnSingleArticle(w http.ResponseWriter, r *http.Request) {
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

// CreateNewArticle controller
func CreateNewArticle(w http.ResponseWriter, r *http.Request) {
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

// DeleteSingleArticle controller
func DeleteSingleArticle(w http.ResponseWriter, r *http.Request) {
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
		models.DB.Delete(&article)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		models.DB.Find(&models.Comments)
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

// UpdateSingleArticle controller
func UpdateSingleArticle(w http.ResponseWriter, r *http.Request) {
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
		models.DB.First(&article, id)
		json.Unmarshal(reqBody, &reqMap)
		if uint(senderId) != article.UserID && sender.Admin == false {
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

// ReturnAllUsers controller
func ReturnAllUsers(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.RawQuery
	// If rawQuery is exists in request. ex: admin=true
	if len(rawQuery) > 0 {
		fmt.Printf("Endpoint Hit: ReturnAllUsers by '%s'\n", rawQuery)
		cut := strings.Split(rawQuery, "=")
		key, value := cut[0], cut[1]
		found := filter(key, value, models.Users)
		if found == nil {
			result := fmt.Sprintf("No user found by '%s': '%s'!", key, value)
			w.WriteHeader(404)
			http.Error(w, result, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
		return
	}
	fmt.Println("Endpoint Hit: returnAllUsers")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Users)
}

// ReturnSingleUser controller
func ReturnSingleUser(w http.ResponseWriter, r *http.Request) {
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

// CreateNewUser controller
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
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
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
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

// DeleteSingleUser controller
func DeleteSingleUser(w http.ResponseWriter, r *http.Request) {
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
		models.DB.Delete(&user)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		models.DB.Find(&models.Comments)
		models.DB.Preload("Comments").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		result := user
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No user found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// UpdateSingleUser controller
func UpdateSingleUser(w http.ResponseWriter, r *http.Request) {
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
		models.DB.First(&user, id)
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
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		result := user
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No user found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)

	}
}

// LoginUser controller
func LoginUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reqMap map[string]string
	json.Unmarshal(reqBody, &reqMap)
	reqUser := reqMap["username"]
	reqPass := reqMap["password"]
	hasher := md5.New()
	hasher.Write([]byte(reqPass))
	hashPass := hex.EncodeToString(hasher.Sum(nil))
	for _, user := range models.Users {
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

// ReturnAllComments controller
func ReturnAllComments(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.RawQuery
	// If rawQuery is exists in request. ex: UserID=1
	if len(rawQuery) > 0 {
		fmt.Printf("Endpoint Hit: ReturnAllComments by '%s'\n", rawQuery)
		cut := strings.Split(rawQuery, "=")
		key, value := cut[0], cut[1]
		found := filter(key, value, models.Comments)
		if found == nil {
			result := fmt.Sprintf("No comment found by '%s': '%s'!", key, value)
			w.WriteHeader(404)
			http.Error(w, result, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
		return
	}
	fmt.Println("Endpoint Hit: returnAllComments")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Comments)
}

// ReturnSingleComment controller
func ReturnSingleComment(w http.ResponseWriter, r *http.Request) {
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

// CreateNewComment controller
func CreateNewComment(w http.ResponseWriter, r *http.Request) {
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
		models.DB.Create(&comment)
		// Reload Users list
		models.DB.Find(&models.Comments)
		models.DB.Preload("Comments").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		json.NewEncoder(w).Encode(comment)
	} else {
		result := fmt.Sprintf("One comment found by id: '%v'!", comment.ID)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// DeleteSingleComment controller
func DeleteSingleComment(w http.ResponseWriter, r *http.Request) {
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
		models.DB.Delete(&comment)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		models.DB.Find(&models.Comments)
		models.DB.Preload("Comments").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		result := comment
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)
	}
}

// UpdateSingleComment controller
func UpdateSingleComment(w http.ResponseWriter, r *http.Request) {
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
		models.DB.First(&comment, id)
		json.Unmarshal(reqBody, &reqMap)
		if senderId != comment.UserID && sender.Admin == false {
			http.Error(w, "Permission Dinied!", http.StatusForbidden)
			return
		}
		comment.Message = reqMap["Message"]
		models.DB.Save(&comment)
		models.DB.Find(&models.Comments)
		models.DB.Preload("Comments").Find(&models.Articles)
		models.DB.Preload("Articles").Preload("Comments").Find(&models.Users)
		models.AppCache.Set("comments", models.Comments, 24*time.Hour)
		models.AppCache.Set("articles", models.Articles, 24*time.Hour)
		models.AppCache.Set("users", models.Users, 24*time.Hour)
		result := comment
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("No comment found by id: '%v'!", id)
		http.Error(w, result, http.StatusBadRequest)

	}
}
