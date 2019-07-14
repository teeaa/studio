package bookings

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
	"github.com/teeaa/studio/internal/classes"
)

func setup() {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	gormDB, _ := gorm.Open(mocket.DriverName, "")
	db = Database{gormDB}
}

// Mock get class by id response for checking existing classes to book to
func setClassMatch() {
	commonReply := []map[string]interface{}{{
		"id":         1,
		"name":       "Class #1",
		"start_date": time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		"end_date":   time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
		"capacity":   20,
	}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "classes"  WHERE ("classes"."id" = 1) ORDER BY "classes"."id" ASC LIMIT 1`).WithReply(commonReply)
}

func TestGetBookingsEmpty(t *testing.T) {
	setup()
	commonReply := []map[string]interface{}{{}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM bookings`).WithReply(commonReply)

	w, r, _ := makeRequest(nil, nil)
	getBookings(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d instead", w.Code)
	}

	var bookings []Booking
	err = json.Unmarshal(body, &bookings)
	if err != nil {
		t.Error("Failed to serialise response to json:", err)
	}
	if len(bookings) != 0 {
		t.Errorf("Expected bookings array length to be 0, was %d instead", len(bookings))
	}
}

func TestAddBooking(t *testing.T) {
	setClassMatch()
	mocket.Catcher.NewMock().WithQuery(`INSERT INTO "bookings"`)

	requestData := Booking{
		123, // Sent ID shouldn't affect result
		"Another Tester",
		time.Date(2019, 8, 11, 0, 0, 0, 0, time.UTC),
		1,
	}

	w, r, _ := makeRequest(&requestData, nil)
	classes.SetupExternally(classes.Database(db))
	addBooking(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusCreated {
		t.Errorf("Expected HTTP status 201 OK, got %d instead", w.Code)
	}

	var booking Booking
	err = json.Unmarshal(body, &booking)
	if err != nil {
		t.Error("Failed to unmarshalling response to json:", err)
	}

	compare := Booking{
		0,
		"Another Tester",
		time.Date(2019, 8, 11, 0, 0, 0, 0, time.UTC),
		1,
	}
	if compare != booking {
		t.Error("Received booking data didn't match expectations:", compare, booking)
	}
}

func TestGetBookingsData(t *testing.T) {
	commonReply := []map[string]interface{}{{
		"id":           1,
		"name":         "Another Tester",
		"booking_date": time.Date(2019, 8, 11, 0, 0, 0, 0, time.UTC),
		"class_id":     1,
	}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "bookings"`).WithReply(commonReply)

	w, r, _ := makeRequest(nil, nil)
	getBookings(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d instead", w.Code)
	}

	compare, _ := json.Marshal([]Booking{{
		1,
		"Another Tester",
		time.Date(2019, 8, 11, 0, 0, 0, 0, time.UTC),
		1,
	}})

	// Need to compare response without Unmarshal because that would reset ids
	compareStr := "'" + strings.TrimSpace(string(compare)) + "'"
	bodyStr := "'" + strings.TrimSpace(string(body)) + "'"
	if compareStr != bodyStr {
		t.Error("Received data didn't match expectations:", compareStr, bodyStr)
	}
}

func TestPutBooking(t *testing.T) {
	commonReply := []map[string]interface{}{{
		"id":           1,
		"name":         "Old tester name",
		"booking_date": time.Date(2019, 8, 11, 0, 0, 0, 0, time.UTC),
		"class_id":     1,
	}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "bookings"  WHERE ("bookings"."id" = 1) ORDER BY "bookings"."id" ASC LIMIT 1`).WithReply(commonReply)
	mocket.Catcher.NewMock().WithQuery(`UPDATE "bookings" SET "name" = ?, "booking_date" = ?, "class_id" = ?  WHERE "bookings"."id" = ?`)

	requestData := Booking{
		123, // Sent ID shouldn't affect result
		"New name",
		time.Date(2019, 8, 15, 0, 0, 0, 0, time.UTC),
		1,
	}

	commonReply[0]["name"] = requestData.Name
	commonReply[0]["booking_date"] = requestData.BookingDate
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "bookings"  WHERE ("bookings"."id" = 1) ORDER BY "bookings"."id" ASC LIMIT 1`).WithReply(commonReply)
	setClassMatch()

	w, r, _ := makeRequest(&requestData, map[string]string{"id": "1"})
	updateBooking(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d instead", w.Code)
	}

	compare := Booking{
		1,
		"New name",
		time.Date(2019, 8, 15, 0, 0, 0, 0, time.UTC),
		1,
	}

	var responseBooking Booking
	json.Unmarshal(body, &responseBooking)
	responseBooking.ID = 1

	// Compare with hard set ID because keys in json structure might change places
	if compare != responseBooking {
		t.Error("Received data didn't match expectations:", compare, responseBooking)
	}
}

func TestDeleteBooking(t *testing.T) {
	w, r, _ := makeRequest(nil, map[string]string{"id": "1"})

	mocket.Catcher.NewMock().WithQuery(`DELETE FROM "boookings"  WHERE "bookings"."id" = ?`)
	deleteBooking(w, r)

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

	if responseBody["message"] != "Booking removed" {
		t.Error("Response body didn't match expectations:", responseBody)
	}
}

func TestGetBooking(t *testing.T) {
	w, r, _ := makeRequest(nil, map[string]string{"id": "1"})

	getBooking(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Error reading body:", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d instead", w.Code)
	}

	compare := Booking{
		0,
		"New name",
		time.Date(2019, 8, 15, 0, 0, 0, 0, time.UTC),
		1,
	}

	var responseBody Booking
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Error("Error unmarshalling body:", err)
	}

	if compare != responseBody {
		t.Error("Response body didn't match expectations:", compare, responseBody)
	}
}

func TestGetBookingNonExisting(t *testing.T) {
	w, r, _ := makeRequest(nil, map[string]string{"id": "123"})

	getBooking(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404, got %d instead", w.Code)
	}
}

func TestPutBookingNonExisting(t *testing.T) {
	requestData := Booking{
		234,
		"Shouldn't work",
		time.Date(2019, 6, 10, 0, 0, 0, 0, time.UTC),
		122,
	}
	w, r, _ := makeRequest(&requestData, map[string]string{"id": "234"})

	updateBooking(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404, got %d instead", w.Code)
	}
}

func TestDeleteBookingNonExisting(t *testing.T) {
	w, r, _ := makeRequest(nil, map[string]string{"id": "345"})

	deleteBooking(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404, got %d instead", w.Code)
	}
}

func TestAddBookingNonExistingClass(t *testing.T) {
	setClassMatch()
	mocket.Catcher.NewMock().WithQuery(`INSERT INTO "bookings"`)

	requestData := Booking{
		123, // Sent ID shouldn't affect result
		"Another Tester",
		time.Date(2019, 8, 11, 0, 0, 0, 0, time.UTC),
		1234,
	}

	w, r, _ := makeRequest(&requestData, nil)
	classes.SetupExternally(classes.Database(db))
	addBooking(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 OK, got %d instead", w.Code)
	}
}

func TestPutBookingNonExistingClass(t *testing.T) {
	setClassMatch()
	mocket.Catcher.NewMock().WithQuery(`INSERT INTO "bookings"`)

	requestData := Booking{
		123, // Sent ID shouldn't affect result
		"Another Tester",
		time.Date(2019, 8, 11, 0, 0, 0, 0, time.UTC),
		221,
	}

	w, r, _ := makeRequest(&requestData, nil)
	classes.SetupExternally(classes.Database(db))
	updateBooking(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 OK, got %d instead", w.Code)
	}
}

func makeRequest(requestData *Booking, vars map[string]string) (*httptest.ResponseRecorder, *http.Request, error) {
	requestBody, err := json.Marshal(&requestData)
	if err != nil {
		return nil, nil, err
	}

	r := httptest.NewRequest("PUT", "/bookings/1", bytes.NewReader(requestBody))
	r.Header.Add("Content-Type", "application/json")
	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()
	w.Header().Add("Content-Type", "application/json")

	return w, r, nil
}
