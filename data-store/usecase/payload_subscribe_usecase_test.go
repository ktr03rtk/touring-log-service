package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
	"github.com/ktr03rtk/touring-log-service/data-store/mock"
	"github.com/stretchr/testify/assert"
)

func TestPayloadSubscribeUsecaseExecute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		returnErr   error
		expectedErr error
	}{
		{
			"normal case",
			nil,
			nil,
		},
		{
			"normal case",
			errors.New("error occurred"),
			errors.New("failed to execute payload subscribe usecase"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mock.NewMockPayloadSubscribeRepository(ctrl)
			usecase := NewPayloadSubscribeUsecase(repository)

			ctx := context.Background()
			ch := make(chan *model.Payload, concurrency)

			repository.EXPECT().Subscribe(ctx, gomock.Any()).Return(tt.returnErr).Times(1)

			if err := usecase.Execute(ctx, ch); err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
			}
		})
	}
}
