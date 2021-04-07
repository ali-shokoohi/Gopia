package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
)

func urlMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request at: %v\n", r.URL)
		handler.ServeHTTP(w, r)
	})
}

func authMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ignore GET method
		if r.Method == "GET" {
			handler.ServeHTTP(w, r)
			return
		}
		reqUser, reqPass, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Access Dinied!", http.StatusForbidden)
			return
		}
		hasher := md5.New()
		hasher.Write([]byte(reqPass))
		hashPass := hex.EncodeToString(hasher.Sum(nil))
		for _, user := range Users {
			if user.Username == reqUser && user.Password == hashPass {
				// Set user in request
				userURL := url.User(fmt.Sprint(user.ID))
				r.URL.User = userURL
				handler.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Access Dinied!", http.StatusForbidden)
	})
}
