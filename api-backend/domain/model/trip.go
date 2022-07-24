package model

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

func NewTrip(id TripID, year, month, day int, unit string) *Trip {
	data := &Trip{
		ID:    id,
		Year:  year,
		Month: month,
		Day:   day,
		Unit:  unit,
	}

	return data
}
