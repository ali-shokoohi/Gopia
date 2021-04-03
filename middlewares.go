package main

import (
	"fmt"
	"net/http"
)

func urlMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request at: %v\n", r.URL)
		handler.ServeHTTP(w, r)
	})
}
