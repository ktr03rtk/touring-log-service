package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTrip(t *testing.T) {
	t.Parallel()

	id := TripID("62c24944-f532-4c5d-a695-70fa3e72f3ab")
	year := 2022
	month := 0o1
	day := 25
	unit := "edge"

	expectedOutput := &Trip{ID: id, Year: year, Month: month, Day: day, Unit: unit}

	output := NewTrip(id, year, month, day, unit)
	assert.Exactly(t, expectedOutput, output)
}
