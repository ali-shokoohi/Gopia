package middlewares

import (
	"fmt"
	"net/http"
)

// URLMiddleWare log requests
func URLMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request at: %v\n", r.URL)
		handler.ServeHTTP(w, r)
	})
}
