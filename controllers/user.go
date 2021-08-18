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

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gitlab.com/greenly/go-rest-api/models"
)

// ReturnAllUsers : Return all users with or without raw query in request
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

// ReturnSingleUser : Return a user article by ID
func ReturnSingleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingeUser by id='%v'\n", id)
	found := findObject(id, models.Users)
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

// CreateNewUser : Create a new user object
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user models.User
	json.Unmarshal(reqBody, &user)
	fmt.Printf("Endpoint Hit: CreateNewUser by id='%v'\n", user.ID)
	found := findObject(fmt.Sprint(user.ID), models.Users)
	if found == nil {
		// Only admins can create super users!
		senderId := r.Context().Value("user")
		if senderId == nil {
			user.Admin = false
		} else {
			senderFound := findObject(fmt.Sprint(senderId.(uint)), models.Users)
			sender := senderFound.(map[string]interface{})
			if sender["Admin"] == false {
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

// DeleteSingleUser : Delete a single user object from database by ID
func DeleteSingleUser(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findObject(fmt.Sprint(senderId), models.Users)
	sender := senderFound.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: deleteSingleUser by id='%v'\n", id)
	// Only owners or admins can do this
	if strconv.Itoa(int(senderId)) != id && sender["Admin"] == false {
		http.Error(w, "Permission Dinied!", http.StatusForbidden)
		return
	}
	found := findObject(id, models.Users)
	if found != nil {
		user := found.(map[string]interface{})
		models.DB.Delete(&models.User{}, user["ID"])
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

// UpdateSingleUser : Update and change a single user in database via ID
func UpdateSingleUser(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("user").(uint)
	senderFound := findObject(fmt.Sprint(senderId), models.Users)
	sender := senderFound.(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Endpoint Hit: updateSingleUser by id='%v'\n", id)
	// Only owners or admins can do this
	if strconv.Itoa(int(senderId)) != id && sender["Admin"] == false {
		http.Error(w, "Permission Dinied!", http.StatusForbidden)
		return
	}
	found := findObject(id, models.Users)
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
		if sender["Admin"] == true {
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

// LoginUser : Return JWT token of entered user if Its username and password are correct
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
