package bookings

import (
	"testing"
	"time"
)

func TestGetBookingFromReq(t *testing.T) {
	w, r, _ := makeRequest(nil, map[string]string{"id": "1"})

	booking, err := db.getBookingFromReq(w, r)
	if err != nil {
		t.Error("Error getting booking by request vars")
	}
	compare := &Booking{
		1,
		"New name",
		time.Date(2019, 8, 15, 0, 0, 0, 0, time.UTC),
		1,
	}
	if *compare != *booking {
		t.Error("Retrieved class data didn't match expectations:", *compare, *booking)
	}
}
