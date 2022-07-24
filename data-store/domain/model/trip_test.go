package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTrip(t *testing.T) {
	t.Parallel()

	id := TripID("62c24944-f532-4c5d-a695-70fa3e72f3ab")
	year := 2022
	month := 0o1
	day := 25
	time := time.Date(2022, 1, 25, 10, 10, 10, 0, time.Local)
	unit := "edge"

	expectedOutput := &Trip{ID: id, Year: year, Month: month, Day: day, Unit: unit}

	output := NewTrip(id, time, unit)
	assert.Exactly(t, expectedOutput, output)
}
