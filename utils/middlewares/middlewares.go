package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/greenly/go-rest-api/models"
)

// CORSMiddleWare for allowing CORS
func CORSMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			fmt.Fprintf(w, "Ok!")
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// JWTMiddleWare for JWT authentication
func JWTMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remove all request user context
		ctx := context.WithValue(r.Context(), "user", nil)
		r = r.WithContext(ctx)
		// Ignore GET method
		if r.Method == "GET" {
			handler.ServeHTTP(w, r)
			return
		}
		notAuth := []string{"/users/new", "/users/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                         //current request path
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
		tk := &models.Token{}

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
		if tk.UserID <= 0 {
			if !necessary[0] {
				handler.ServeHTTP(w, r)
				return
			}
			http.Error(w, "User is not valid anymore!", http.StatusForbidden)
			return
		}
		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Printf("User %v:\n", tk.UserID) //Useful for monitoring
		ctx = context.WithValue(r.Context(), "user", tk.UserID)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
