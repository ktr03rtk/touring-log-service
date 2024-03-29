package model

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPayload(t *testing.T) {
	t.Parallel()

	message := []byte("test message")
	tests := []struct {
		name           string
		key            string
		expectedOutput *Payload
		expectedErr    error
	}{
		{
			"normal case",
			"touring-log/raw/thing=thingName/year=2022/month=01/day=12/2022-01-12-12-51-10.dat",
			&Payload{message: message, key: "touring-log/raw/thing=thingName/year=2022/month=01/day=12/2022-01-12-12-51-10.dat"},
			nil,
		},
		{
			"validate error case: month",
			"touring-log/raw/thing=thingName/year=2022/month=21/day=12/2022-01-12-12-51-10.dat",
			nil,
			errors.New("failed to validate key"),
		},
		{
			"validate error case: object",
			"touring-log/raw/thing=thingName/year=2022/month=01/day=12/2022-01-12-12-error.dat",
			nil,
			errors.New("failed to validate key"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := NewPayload(message, tt.key)
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

func TestGetDate(t *testing.T) {
	t.Parallel()

	const shortForm = "2006-01-02"
	date := time.Date(2022, 1, 12, 0, 0, 0, 0, time.UTC)

	message := []byte("test message")
	tests := []struct {
		name           string
		input          Payload
		expectedOutput *time.Time
		expectedErr    error
	}{
		{
			"normal case",
			Payload{message: message, key: "touring-log/raw/thing=thingName/year=2022/month=01/day=12/2022-01-12-12-51-10.dat"},
			&date,
			nil,
		},
		{
			"validate error case: month",
			Payload{message: message, key: "touring-log/raw/thing=thingName/year=2022/month=21/day=12/2022-01-12-12-51-10.dat"},
			nil,
			errors.New("failed to parse date"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := tt.input.GetDate()
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

func TestGetUnit(t *testing.T) {
	t.Parallel()

	message := []byte("test message")
	payload := Payload{message: message, key: "touring-log/raw/thing=thingName/year=2022/month=01/day=12/2022-01-12-12-51-10.dat"}
	expectedOutput := "thingName"

	output := payload.GetUnit()
	assert.Exactly(t, expectedOutput, output)
}
