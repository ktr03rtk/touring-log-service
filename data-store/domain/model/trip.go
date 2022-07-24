package model

import (
	"time"
)

type Trip struct {
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Day   int    `json:"day"`
	Unit  string `json:"unit"`
}

func NewTrip(date time.Time, unit string) *Trip {
	data := &Trip{
		Year:  date.Year(),
		Month: int(date.Month()),
		Day:   date.Day(),
		Unit:  unit,
	}

	return data
}
