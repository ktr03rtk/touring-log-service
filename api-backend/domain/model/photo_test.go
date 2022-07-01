package model

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPhoto(t *testing.T) {
	t.Parallel()

	id := PhotoID("62c24944-f532-4c5d-a695-70fa3e72f3ab")
	userID := UserID("72c24944-f532-4c5d-a695-70fa3e72f3ab")
	year := 2022
	month := 0o1
	day := 25
	time := time.Date(2022, 1, 25, 10, 10, 10, 0, time.Local)
	unit := "edge"

	tests := []struct {
		name           string
		lat            float64
		lon            float64
		expectedOutput *Photo
		expectedErr    error
	}{
		{
			"normal case",
			35.470403,
			139.625228,
			&Photo{ID: id, Year: year, Month: month, Day: day, Lat: 35.470403, Lon: 139.625228, Timestamp: time, S3ObjectKey: "touring-log/photo/thing=edge/year=2022/month=01/day=25/1643073010000.gz", UserID: userID},
			nil,
		},
		{
			"normal case: negative gps value",
			-35.470403,
			-139.625228,
			&Photo{ID: id, Year: year, Month: month, Day: day, Lat: -35.470403, Lon: -139.625228, Timestamp: time, S3ObjectKey: "touring-log/photo/thing=edge/year=2022/month=01/day=25/1643073010000.gz", UserID: userID},
			nil,
		},
		{
			"validate error case: lat",
			95.470403,
			139.625228,
			nil,
			errors.New("failed to satisfy GPS Lat Spec"),
		},
		{
			"validate error case: lat",
			35.470403,
			239.625228,
			nil,
			errors.New("failed to satisfy GPS Lon Spec"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := NewPhoto(id, time, tt.lat, tt.lon, userID, unit)
			if err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
				assert.Exactly(t, tt.expectedOutput, output)
			}
		})
	}
}
