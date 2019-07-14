package bookings

import (
	"encoding/json"
	"errors"
	"time"
)

// Booking representation of bookings.bookings
type Booking struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	BookingDate time.Time `gorm:"type:date" json:"booking_date"`
	ClassID     uint64    `json:"class_id"`
}

// MarshalJSON to date correctly
func (b *Booking) MarshalJSON() ([]byte, error) {
	type Alias Booking
	return json.Marshal(&struct {
		BookingDate string `json:"booking_date"`
		*Alias
	}{
		BookingDate: b.BookingDate.Format("2006-01-02"),
		Alias:       (*Alias)(b),
	})
}

// UnmarshalJSON to date correctly and strip ID field from requests
func (b *Booking) UnmarshalJSON(data []byte) error {
	type Alias Booking
	aux := &struct {
		ID          uint64 `gorm:"-" sql:"-" json:"id"`
		BookingDate string `json:"booking_date"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if len(aux.BookingDate) < 10 {
		return errors.New("Invalid booking_date in payload")
	}

	b.BookingDate, err = time.Parse("2006-01-02", aux.BookingDate[0:10])
	if err != nil {
		return err
	}

	return nil
}
