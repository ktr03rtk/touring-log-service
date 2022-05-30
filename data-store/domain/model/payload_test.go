package model

import (
	"errors"
	"testing"

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
			"touring-log/raw/thing=thingName/month=01/day=12/2022-01-12-12-51.dat",
			&Payload{message: message, key: "touring-log/raw/thing=thingName/month=01/day=12/2022-01-12-12-51.dat"},
			nil,
		},
		{
			"validate error case: month",
			"touring-log/raw/thing=thingName/month=21/day=12/2022-01-12-12-51.dat",
			nil,
			errors.New("failed to validate key"),
		},
		{
			"validate error case: object",
			"touring-log/raw/thing=thingName/month=01/day=12/2022-01-12-12-error.dat",
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
