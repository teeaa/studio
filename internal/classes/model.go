package classes

import (
	"encoding/json"
	"errors"
	"time"
)

// Class representation of classes.classes
type Class struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `gorm:"type:date" json:"start_date"`
	EndDate   time.Time `gorm:"type:date" json:"end_date"`
	Capacity  uint      `json:"capacity"`
}

// MarshalJSON to date correctly
func (c *Class) MarshalJSON() ([]byte, error) {
	type Alias Class
	return json.Marshal(&struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		*Alias
	}{
		StartDate: c.StartDate.Format("2006-01-02"),
		EndDate:   c.EndDate.Format("2006-01-02"),
		Alias:     (*Alias)(c),
	})
}

// UnmarshalJSON to date correctly and strip ID field from requests
func (c *Class) UnmarshalJSON(data []byte) error {
	type Alias Class
	aux := &struct {
		ID        uint64 `gorm:"-" sql:"-" json:"id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if aux.Capacity < 1 {
		return errors.New("Invalid capacity in payload, must to be over 0")
	}

	if aux.StartDate != "" {
		if len(aux.StartDate) < 10 {
			return errors.New("Invalid start_date in payload")
		}

		c.StartDate, err = time.Parse("2006-01-02", aux.StartDate[0:10])
		if err != nil {
			return err
		}
	}

	if aux.EndDate != "" {
		if len(aux.EndDate) < 10 {
			return errors.New("Invalid end_date in payload")
		}

		c.EndDate, err = time.Parse("2006-01-02", aux.EndDate[0:10])
		if err != nil {
			return err
		}
	}

	return nil
}
