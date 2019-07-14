package bookings

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/teeaa/studio/internal/classes"
	"github.com/teeaa/studio/internal/helpers"
)

func getBookings(w http.ResponseWriter, r *http.Request) {
	var bookings []Booking
	err := db.Find(&bookings).Error

	if err != nil {
		log.Error("Error fetching bookings from db: ", err)
		helpers.ResponseJSON(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	json.NewEncoder(w).Encode(&bookings)
}

func checkValidity(booking Booking, class classes.Class) error {
	if class == (classes.Class{}) {
		return errors.New("No such class")
	}

	class.EndDate = class.EndDate.Add(24 * time.Hour)
	booking.BookingDate = booking.BookingDate.Add(12 * time.Hour)

	if !booking.BookingDate.After(class.StartDate) || !booking.BookingDate.Before(class.EndDate) {
		return errors.New("Booked date outside class start and end")
	}

	return nil
}

func addBooking(w http.ResponseWriter, r *http.Request) {
	var booking Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		log.Warn("Error parsing JSON when creating new booking: ", err)
		helpers.ResponseJSON(w, http.StatusBadRequest, "Invalid request body for booking")
		return
	}

	class, err := classes.GetClassByID(booking.ClassID)
	if err != nil {
		log.Warn("Tried to book with non-existing class id")
		helpers.ResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = checkValidity(booking, class)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.Create(&booking).Error
	if err != nil {
		log.Error("Error inserting booking to db: ", err)
		helpers.ResponseJSON(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&booking)
}

func getBooking(w http.ResponseWriter, r *http.Request) {
	booking, err := db.getBookingFromReq(w, r)
	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(&booking)
}

func updateBooking(w http.ResponseWriter, r *http.Request) {
	booking, err := db.getBookingFromReq(w, r)
	if err != nil {
		return
	}

	err = json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		log.Warn("Error parsing JSON when creating new booking: ", err)
		helpers.ResponseJSON(w, http.StatusBadRequest, "Invalid request body for booking")
		return
	}

	class, err := classes.GetClassByID(booking.ClassID)

	err = checkValidity(*booking, class)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.Save(&booking).Error
	if err != nil {
		log.Error("Error saving booking to db: ", err)
		helpers.ResponseJSON(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	json.NewEncoder(w).Encode(&booking)
}
func deleteBooking(w http.ResponseWriter, r *http.Request) {
	booking, err := db.getBookingFromReq(w, r)
	if err != nil {
		return
	}

	err = db.Delete(&booking).Error
	if err != nil {
		log.Error("Error deleting booking from db: ", err)
		helpers.ResponseJSON(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	helpers.ResponseJSON(w, 200, "Booking removed")
}

// Routes set routes for /bookings
func Routes(gormDB *gorm.DB, router *mux.Router) {
	db = Database{gormDB}
	db.AutoMigrate(&Booking{})

	router.HandleFunc("", getBookings).Methods("GET")
	router.HandleFunc("", addBooking).Methods("POST")
	router.HandleFunc("/{id}", getBooking).Methods("GET")
	router.HandleFunc("/{id}", updateBooking).Methods("PUT")
	router.HandleFunc("/{id}", deleteBooking).Methods("DELETE")
}
