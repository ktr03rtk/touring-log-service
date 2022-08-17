package usecase

import (
	"testing"

	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/mock"
	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDateListQueryUsecaseCase(t *testing.T) {
	t.Parallel()

	user_id := "72c24944-f532-4c5d-a695-70fa3e72f3ab"
	unit := "edge"
	query := "SELECT year, month, day FROM photos WHERE year = ? AND month = ? AND user_id = ? UNION SELECT year, month, day FROM trips WHERE year = ? AND month = ? AND unit = ?"

	tests := []struct {
		name           string
		year           int
		month          int
		day            int
		fetchOutput    []*model.LogDate
		fetchErr       error
		expectedOutput []*model.LogDate
		expectedErr    error
	}{
		{
			"normal case",
			2022,
			8,
			11,
			[]*model.LogDate{{Year: 8, Month: 8, Day: 11}},
			nil,
			[]*model.LogDate{{Year: 8, Month: 8, Day: 11}},
			nil,
		},
		{
			"fetch error case",
			2022,
			8,
			11,
			nil,
			errors.New("failed to fetch"),
			nil,
			errors.New("failed to execute list query usecase"),
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			queryAdapterRepository := mock.NewMockQueryRepository(ctrl)
			usecase := NewDateListQueryUsecase(queryAdapterRepository)

			args := []interface{}{
				tt.year, tt.month, user_id, tt.year, tt.month, unit,
			}

			var touringLog []*model.LogDate

			gomock.InOrder(
				queryAdapterRepository.EXPECT().Fetch(query, args, touringLog).Return(tt.fetchOutput, tt.fetchErr).Times(1),
			)

			output, err := usecase.Execute(tt.year, tt.month, user_id, unit)
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
