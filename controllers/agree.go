package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/greenly/go-rest-api/models"
)

// ReturnAllAgrees - Return all agrees with or without raw query in request
func ReturnAllAgrees(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.RawQuery
	// If rawQuery is exists in request. ex: UserID=1
	if len(rawQuery) > 0 {
		fmt.Printf("Endpoint Hit: ReturnAllAgrees by '%s'\n", rawQuery)
		cut := strings.Split(rawQuery, "=")
		key, value := cut[0], cut[1]
		found := filter(key, value, models.Agrees)
		if found == nil {
			result := fmt.Sprintf("No agree found by '%s': '%s'!", key, value)
			w.WriteHeader(404)
			http.Error(w, result, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
		return
	}
	fmt.Println("Endpoint Hit: returnAllAgrees")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Agrees)
}
