package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func urlMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request at: %v\n", r.URL)
		handler.ServeHTTP(w, r)
	})
}

// func authMiddleWare(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Ignore GET method
// 		if r.Method == "GET" {
// 			handler.ServeHTTP(w, r)
// 			return
// 		}
// 		reqUser, reqPass, ok := r.BasicAuth()
// 		if !ok {
// 			http.Error(w, "Access Dinied!", http.StatusForbidden)
// 			return
// 		}
// 		hasher := md5.New()
// 		hasher.Write([]byte(reqPass))
// 		hashPass := hex.EncodeToString(hasher.Sum(nil))
// 		for _, user := range Users {
// 			if user.Username == reqUser && user.Password == hashPass {
// 				// Set user in request
// 				userURL := url.User(fmt.Sprint(user.ID))
// 				r.URL.User = userURL
// 				handler.ServeHTTP(w, r)
// 				return
// 			}
// 		}
// 		http.Error(w, "Access Dinied!", http.StatusForbidden)
// 	})
// }

func CORSMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		handler.ServeHTTP(w, r)
	})
}

func jwtMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remove all request user context
		ctx := context.WithValue(r.Context(), "user", nil)
		r = r.WithContext(ctx)
		// Ignore GET method
		if r.Method == "GET" {
			handler.ServeHTTP(w, r)
			return
		}
		notAuth := []string{"/user/new", "/user/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                       //current request path
		necessary := [1]bool{true}
		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				necessary[0] = false
				break
			}
		}
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
		if tokenHeader == "" {
			if !necessary[0] {
				handler.ServeHTTP(w, r)
				return
			} //Token is missing, returns with error code 403 Unauthorized
			http.Error(w, "Missing auth token", http.StatusForbidden)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			if !necessary[0] {
				handler.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Invalid/Malformed auth token", http.StatusForbidden)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
		if err != nil { //Malformed token, returns with http code 403 as usual
			if !necessary[0] {
				handler.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Malformed authentication token", http.StatusForbidden)
			return
		}
		if !token.Valid { //Token is invalid, maybe not signed on this server
			if !necessary[0] {
				handler.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Token is not valid.", http.StatusForbidden)
			return
		}
		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Printf("User %v", tk.UserId) //Useful for monitoring
		ctx = context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
