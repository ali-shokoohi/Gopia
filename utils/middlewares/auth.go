package middlewares

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"gitlab.com/greenly/go-rest-api/models"
)

// AuthMiddleWare : Basic authentication middleware
func AuthMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ignore GET method
		if r.Method == "GET" {
			handler.ServeHTTP(w, r)
			return
		}
		notAuth := []string{"/users/new", "/users/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                         //current request path
		necessary := true
		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				necessary = false
				break
			}
		}
		reqUser, reqPass, ok := r.BasicAuth()
		if !ok {
			if necessary {
				handler.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Access Dinied!", http.StatusForbidden)
			return
		}
		hasher := md5.New()
		hasher.Write([]byte(reqPass))
		hashPass := hex.EncodeToString(hasher.Sum(nil))
		for _, user := range models.Users {
			if user.Username == reqUser && user.Password == hashPass {
				// Set user in request
				ctx := context.WithValue(r.Context(), "user", user.ID)
				r = r.WithContext(ctx)
				handler.ServeHTTP(w, r)
				return
			}
		}
		if !necessary {
			handler.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Access Dinied!", http.StatusForbidden)
	})
}
