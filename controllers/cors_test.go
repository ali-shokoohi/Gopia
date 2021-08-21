package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/greenly/go-rest-api/routers"
)

// Test for SkipCORS controller
func TestSkipCORS(t *testing.T) {
	// Get routter object
	routter := routers.CreateRouter()
	// Get listening port from environments
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	// Create http requests for each endpoint
	endPoints := []string{
		"", "articles", "users", "users/login",
		"users/new", "comments", "likes", "agrees",
	}
	for _, endPoint := range endPoints {
		request, err := http.NewRequest("OPTIONS", "http://127.0.0.1:"+port+"/"+endPoint, nil)
		response := httptest.NewRecorder()
		routter.ServeHTTP(response, request)
		// Check we have no error
		assert.Nil(t, err, "Request error must be Nil!")
		// Check response status code must be 200
		assert.Equal(t, 200, response.Result().StatusCode, endPoint+"'s Response code be: "+
			fmt.Sprint(200)+"! But It's: "+fmt.Sprint(response.Result().StatusCode)+"!!!")

		expected := "Ok!"
		result := response.Body.String()
		// Check the response body must be equal to expected value
		assert.Equal(t, expected, response.Body.String(), endPoint+"'s Response must be: "+
			expected+"! But It's: "+result+"!!!")
	}
}
