package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
	"github.com/ktr03rtk/touring-log-service/data-store/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestPayloadStoreUsecaseExecute(t *testing.T) {
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
			"error case",
			errors.New("error occurred"),
			errors.New("failed to execute payload store usecase"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mock.NewMockPayloadStoreRepository(ctrl)
			usecase := NewPayloadStoreUsecase(repository)

			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.WithContext(ctx)
			ch := make(chan *model.Payload)

			payload := &model.Payload{}
			repository.EXPECT().Store(ctx, payload).Return(tt.returnErr).Times(1)

			eg.Go(func() error { return usecase.Execute(ctx, ch) })

			ch <- payload
			cancel()

			if err := eg.Wait(); err != nil {
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
