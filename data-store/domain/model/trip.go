package model

import (
	"time"
)

type Trip struct {
	ID    TripID `json:"id"`
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Day   int    `json:"day"`
	Unit  string `json:"unit"`
}

type (
	TripID string
)

func NewTrip(id TripID, time time.Time, unit string) *Trip {
	data := &Trip{
		ID:    id,
		Year:  time.Year(),
		Month: int(time.Month()),
		Day:   time.Day(),
		Unit:  unit,
	}

	return data
}
