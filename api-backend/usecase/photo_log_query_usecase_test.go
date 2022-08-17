package usecase

import (
	"testing"

	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/mock"
	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPhotoStoreUseCase(t *testing.T) {
	t.Parallel()

	id := "72c24944-f532-4c5d-a695-70fa3e72f3ab"

	tests := []struct {
		name           string
		year           int
		month          int
		day            int
		query          string
		fetchOutput    interface{}
		fetchErr       error
		expectedOutput []*model.WebClientPhoto
		expectedErr    error
	}{
		{
			"normal case",
			2022,
			8,
			11,
			"SELECT id, lat, lon as lng FROM photos WHERE year = ? AND month = ? AND day = ? AND user_id = ?",
			[]*model.WebClientPhoto{{Id: "10022073-00ac-493f-9f3c-683725110408", Lat: 35.470403, Lng: 139.625228}, {Id: "10022073-00ac-493f-9f3c-683725110409", Lat: 35.470404, Lng: 139.625228}},
			nil,
			[]*model.WebClientPhoto{{Id: "10022073-00ac-493f-9f3c-683725110408", Lat: 35.470403, Lng: 139.625228}, {Id: "10022073-00ac-493f-9f3c-683725110409", Lat: 35.470404, Lng: 139.625228}},
			nil,
		},
		{
			"error case",
			2022,
			8,
			11,
			"SELECT id, lat, lon as lng FROM photos WHERE year = ? AND month = ? AND day = ? AND user_id = ?",
			nil,
			errors.New("failed to fetch"),
			nil,
			errors.New("failed to execute photo log query usecase"),
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			queryAdapterRepository := mock.NewMockQueryRepository(ctrl)
			usecase := NewPhotoLogQueryUsecase(queryAdapterRepository)

			args := []interface{}{
				tt.year, tt.month, tt.day, id,
			}

			var touringLog []*model.WebClientPhoto

			gomock.InOrder(
				queryAdapterRepository.EXPECT().Fetch(tt.query, args, touringLog).Return(tt.fetchOutput, tt.fetchErr).Times(1),
			)

			output, err := usecase.Execute(tt.year, tt.month, tt.day, id)
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
