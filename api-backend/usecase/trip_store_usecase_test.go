package usecase

import (
	"testing"

	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/service"
	"github.com/ktr03rtk/touring-log-service/api-backend/mock"
	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTripStoreUseCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                    string
		year                    int
		month                   int
		day                     int
		unit                    string
		findByDateAndUnitOutput *model.Trip
		findByDateAndUnitErr    error
		createErr               error
		expectedOutput          error
		expectedCallTimes       int
	}{
		{
			"normal case",
			2022,
			8,
			11,
			"edge",
			nil,
			nil,
			nil,
			nil,
			1,
		},
		{
			"already exsits case",
			2022,
			8,
			11,
			"edge",
			&model.Trip{},
			nil,
			nil,
			nil,
			0,
		},
		{
			"find error case",
			2022,
			8,
			11,
			"edge",
			nil,
			errors.New("failed to find"),
			nil,
			errors.New("failed to find trip"),
			0,
		},
		{
			"create error case",
			2022,
			8,
			11,
			"edge",
			nil,
			nil,
			errors.New("failed to create"),
			errors.New("failed to execute trip store usecase"),
			1,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tripMetadataRepository := mock.NewMockTripMetadataStoreRepository(ctrl)
			tripService := service.NewTripService(tripMetadataRepository)
			usecase := NewTripStoreUsecase(tripMetadataRepository, tripService)

			gomock.InOrder(
				tripMetadataRepository.EXPECT().FindByDateAndUnit(tt.year, tt.month, tt.day, tt.unit).Return(tt.findByDateAndUnitOutput, tt.findByDateAndUnitErr).Times(1),
				tripMetadataRepository.EXPECT().Create(gomock.Any()).Return(tt.createErr).Times(tt.expectedCallTimes),
			)

			if err := usecase.Execute(tt.year, tt.month, tt.day, tt.unit); err != nil {
				if tt.expectedOutput != nil {
					assert.Contains(t, err.Error(), tt.expectedOutput.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Nil(t, tt.expectedOutput, "error is expected but received nil")
			}
		})
	}
}
