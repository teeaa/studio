package classes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
)

func setup() {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	gormDB, _ := gorm.Open(mocket.DriverName, "")
	db = Database{gormDB}
}

func TestGetClassesEmpty(t *testing.T) {
	setup()
	commonReply := []map[string]interface{}{{}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM classes`).WithReply(commonReply)

	w, r, _ := makeRequest(nil, nil)
	getClasses(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d instead", w.Code)
	}

	var classes []Class
	err = json.Unmarshal(body, &classes)
	if err != nil {
		t.Error("Failed to serialise response to json:", err)
	}
	if len(classes) != 0 {
		t.Errorf("Expected classes array length to be 0, was %d instead", len(classes))
	}
}

func TestAddClass(t *testing.T) {
	mocket.Catcher.Reset().NewMock().WithQuery(`INSERT INTO "classes"`)

	requestData := Class{
		123, // Sent ID shouldn't affect result
		"Class #1",
		time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
		20,
	}

	w, r, _ := makeRequest(&requestData, nil)
	addClass(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusCreated {
		t.Errorf("Expected HTTP status 201 OK, got %d instead", w.Code)
	}

	var class Class
	err = json.Unmarshal(body, &class)
	if err != nil {
		t.Error("Failed to unmarshalling response to json:", err)
	}

	compare := Class{
		0,
		"Class #1",
		time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
		20,
	}
	if compare != class {
		t.Error("Received class data didn't match expectations:", compare, class)
	}
}

func TestGetClassesData(t *testing.T) {
	commonReply := []map[string]interface{}{{
		"id":         1,
		"name":       "Class #1",
		"start_date": time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		"end_date":   time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
		"capacity":   20,
	}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "classes"`).WithReply(commonReply)

	w, r, _ := makeRequest(nil, nil)
	getClasses(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d instead", w.Code)
	}

	compare, _ := json.Marshal([]Class{{
		1,
		"Class #1",
		time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
		20,
	}})

	// Need to compare response without Unmarshal because that would reset ids
	compareStr := "'" + strings.TrimSpace(string(compare)) + "'"
	bodyStr := "'" + strings.TrimSpace(string(body)) + "'"
	if compareStr != bodyStr {
		t.Error("Received data didn't match expectations:", compareStr, bodyStr)
	}
}

func TestPutClass(t *testing.T) {
	commonReply := []map[string]interface{}{{
		"id":         1,
		"name":       "Class #1",
		"start_date": time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		"end_date":   time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
		"capacity":   20,
	}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "classes"  WHERE ("classes"."id" = 1) ORDER BY "classes"."id" ASC LIMIT 1`).WithReply(commonReply)
	mocket.Catcher.NewMock().WithQuery(`UPDATE "classes" SET "name" = ?, "start_date" = ?, "end_date" = ?, "capacity" = ?  WHERE "classes"."id" = ?`)

	requestData := Class{
		123, // Sent ID shouldn't affect result
		"New class name",
		time.Date(2019, 6, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 8, 22, 0, 0, 0, 0, time.UTC),
		15,
	}

	commonReply[0]["name"] = requestData.Name
	commonReply[0]["start_date"] = requestData.StartDate
	commonReply[0]["end_date"] = requestData.EndDate
	commonReply[0]["capacity"] = requestData.Capacity
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "classes"  WHERE ("classes"."id" = 1) ORDER BY "classes"."id" ASC LIMIT 1`).WithReply(commonReply)

	w, r, _ := makeRequest(&requestData, map[string]string{"id": "1"})
	updateClass(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d instead", w.Code)
	}

	compare := Class{
		1,
		"New class name",
		time.Date(2019, 6, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 8, 22, 0, 0, 0, 0, time.UTC),
		15,
	}

	var responseClass Class
	json.Unmarshal(body, &responseClass)
	responseClass.ID = 1

	// Compare with hard set ID because keys in json structure might change places
	if compare != responseClass {
		t.Error("Received data didn't match expectations:", compare, responseClass)
	}
}

func TestDeleteClass(t *testing.T) {
	w, r, _ := makeRequest(nil, map[string]string{"id": "1"})

	mocket.Catcher.NewMock().WithQuery(`DELETE FROM "classes"  WHERE "classes"."id" = ?`)
	deleteClass(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d instead", w.Code)
	}

	var responseBody map[string]string
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Error("Error unmarshalling body:", err)
	}

	if responseBody["message"] != "Class removed" {
		t.Error("Response body didn't match expectations:", responseBody)
	}
}

func TestGetClass(t *testing.T) {
	w, r, _ := makeRequest(nil, map[string]string{"id": "1"})

	getClass(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d instead", w.Code)
	}

	compare := Class{
		0,
		"New class name",
		time.Date(2019, 6, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 8, 22, 0, 0, 0, 0, time.UTC),
		15,
	}

	var responseBody Class
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Error("Error unmarshalling body:", err)
	}

	if compare != responseBody {
		t.Error("Response body didn't match expectations:", compare, responseBody)
	}
}

func TestGetClassNonExisting(t *testing.T) {
	w, r, _ := makeRequest(nil, map[string]string{"id": "123"})

	getClass(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404, got %d instead", w.Code)
	}
}

func TestPutClassNonExisting(t *testing.T) {
	requestData := Class{
		234,
		"Shouldn't work",
		time.Date(2019, 6, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 8, 22, 0, 0, 0, 0, time.UTC),
		15,
	}
	w, r, _ := makeRequest(&requestData, map[string]string{"id": "234"})

	updateClass(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404, got %d instead", w.Code)
	}
}

func TestDeleteClassNonExisting(t *testing.T) {
	w, r, _ := makeRequest(nil, map[string]string{"id": "345"})

	deleteClass(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404, got %d instead", w.Code)
	}
}

func makeRequest(requestData *Class, vars map[string]string) (*httptest.ResponseRecorder, *http.Request, error) {
	requestBody, err := json.Marshal(&requestData)
	if err != nil {
		return nil, nil, err
	}

	r := httptest.NewRequest("PUT", "/classes/1", bytes.NewReader(requestBody))
	r.Header.Add("Content-Type", "application/json")
	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()
	w.Header().Add("Content-Type", "application/json")

	return w, r, nil
}
