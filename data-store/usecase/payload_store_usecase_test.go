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

	testPayload, err := model.NewPayload([]byte("test message"), "touring-log/raw/thing=thingName/year=2022/month=01/day=12/2022-01-12-12-51-10.dat")
	assert.Nil(t, err)

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

			str := mock.NewMockPayloadStoreRepository(ctrl)
			tms := mock.NewMockTripMetadataRepository(ctrl)
			usecase := NewPayloadStoreUsecase(str, tms)

			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.WithContext(ctx)
			ch := make(chan *model.Payload)

			str.EXPECT().Store(ctx, testPayload).Return(tt.returnErr).Times(1)
			// TODO: add test case
			tms.EXPECT().Create(gomock.Any()).Return(nil).AnyTimes()

			eg.Go(func() error { return usecase.Execute(ctx, ch) })

			ch <- testPayload
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
