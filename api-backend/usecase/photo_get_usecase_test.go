package usecase

import (
	"io"
	"strings"
	"testing"

	"github.com/ktr03rtk/touring-log-service/api-backend/mock"
	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPhotoGetUseCase(t *testing.T) {
	t.Parallel()

	photo_id := "72c24944-f532-4c5d-a695-70fa3e72f3ab"

	tests := []struct {
		name              string
		query             string
		fetchOutput       []*photoKey
		fetchErr          error
		getOutput         io.ReadCloser
		getErr            error
		expectedOutput    io.ReadCloser
		expectedErr       error
		expectedCallTimes int
	}{
		{
			"normal case",
			"SELECT * FROM photos WHERE id = ?",
			[]*photoKey{{S3ObjectKey: "foo/bar/baz.jpeg.gz"}},
			nil,
			io.NopCloser(strings.NewReader("test")),
			nil,
			io.NopCloser(strings.NewReader("test")),
			nil,
			1,
		},
		{
			"fetch error case",
			"SELECT * FROM photos WHERE id = ?",
			[]*photoKey{{S3ObjectKey: "foo/bar/baz.jpeg.gz"}},
			errors.New("failed to fetch"),
			nil,
			nil,
			nil,
			errors.New("failed to execute photo get query usecase"),
			0,
		},
		{
			"get error case",
			"SELECT * FROM photos WHERE id = ?",
			[]*photoKey{{S3ObjectKey: "foo/bar/baz.jpeg.gz"}},
			nil,
			nil,
			errors.New("failed to get"),
			nil,
			errors.New("failed to get object"),
			1,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			queryAdapterRepository := mock.NewMockQueryRepository(ctrl)
			photoImageRepository := mock.NewMockPhotoImageRepository(ctrl)
			usecase := NewPhotoGetUsecase(queryAdapterRepository, photoImageRepository)

			args := []interface{}{
				photo_id,
			}

			var photoMetadata []*photoKey

			gomock.InOrder(
				queryAdapterRepository.EXPECT().Fetch(tt.query, args, photoMetadata).Return(tt.fetchOutput, tt.fetchErr).Times(1),
				photoImageRepository.EXPECT().Get(tt.fetchOutput[0].S3ObjectKey).Return(tt.getOutput, tt.getErr).Times(tt.expectedCallTimes),
			)

			output, err := usecase.Execute(photo_id)
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
