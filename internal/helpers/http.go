package helpers

import (
	"encoding/json"
	"net/http"
)

/***
 * HTTP helpers
 ***/

// ResponseJSON simple wrapper for sending a JSON response
func ResponseJSON(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
