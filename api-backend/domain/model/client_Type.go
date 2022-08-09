package model

type LogDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type TouringLog struct {
	Trip  []*WebClientTrip
	Photo []*WebClientPhoto
}

type WebClientTrip struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type WebClientPhoto struct {
	Id  string  `json:"id"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
