package bookings

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/teeaa/studio/internal/helpers"
)

// Database wrapper
type Database struct {
	*gorm.DB
}

var db Database

// Get booking from database by id in request and handle error situations
func (db *Database) getBookingFromReq(w http.ResponseWriter, r *http.Request) (*Booking, error) {
	var booking Booking
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Warnf("Requested booking id (%s) is not an integer: %s", vars["id"], err)
		helpers.ResponseJSON(w, http.StatusBadRequest, "Invalid booking ID")
		return nil, err
	}

	err = db.First(&booking, bookingID).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Warnf("Requested booking by id %d does not exist", bookingID)
			helpers.ResponseJSON(w, http.StatusNotFound, "Booking does not exist")
		} else {
			log.Error("Error fetching booking from db: ", err)
			helpers.ResponseJSON(w, http.StatusInternalServerError, "Something went wrong")
		}
		return nil, err
	}
	return &booking, nil
}
