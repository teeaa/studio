package helpers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestResponseJSON(t *testing.T) {
	r := httptest.NewRequest("PUT", "/classes/1", nil)
	r.Header.Add("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	w.Header().Add("Content-Type", "application/json")

	ResponseJSON(w, http.StatusTeapot, "Sent message")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusTeapot {
		t.Errorf("Expected HTTP status 200 OK, got %d instead", w.Code)
	}

	var responseBody map[string]string
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Error("Error unmarshalling body:", err)
	}

	if responseBody["message"] != "Sent message" {
		t.Errorf("Was expecting 'Sent message', got %s instead", responseBody["message"])
	}
}
