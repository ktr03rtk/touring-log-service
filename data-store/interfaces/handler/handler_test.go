package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
	"github.com/ktr03rtk/touring-log-service/data-store/mock"
	"github.com/ktr03rtk/touring-log-service/data-store/usecase"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestPayloadHandler(t *testing.T) {
	// t.Parallel()

	testPayload, err := model.NewPayload([]byte("test message"), "touring-log/raw/thing=thingName/month=01/day=12/2022-01-12-12-51.dat")
	assert.Nil(t, err)

	tests := []struct {
		name               string
		subscribeErr       error
		storeErr           error
		expectedErr        error
		subscribeCallTimes int
		storeCallTimes     int
	}{
		{
			"normal case",
			nil,
			nil,
			nil,
			1,
			1,
		},
		{
			"subscribe error case",
			errors.New("error occurred"),
			nil,
			errors.New("failed to execute payload subscribe usecase"),
			1,
			0,
		},
		{
			"store error case",
			nil,
			errors.New("error occurred"),
			errors.New("failed to execute payload store usecase"),
			1,
			1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sbr := mock.NewMockPayloadSubscribeRepository(ctrl)
			str := mock.NewMockPayloadStoreRepository(ctrl)
			sbu := usecase.NewPayloadSubscribeUsecase(sbr)
			stu := usecase.NewPayloadStoreUsecase(str)

			ph := NewPayloadHandler(sbu, stu)

			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.WithContext(ctx)

			sbr.EXPECT().Subscribe(ctx, gomock.Any()).Return(tt.subscribeErr).Times(tt.subscribeCallTimes)
			str.EXPECT().Store(ctx, testPayload).Return(tt.storeErr).Times(tt.storeCallTimes)

			eg.Go(func() error { return ph.Handle(ctx) })

			h, ok := ph.(*payloadHandler)
			assert.True(t, ok)

			h.payloadCh <- testPayload
			gosched()
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

// Sleep momentarily so that other goroutines can process.
func gosched() { time.Sleep(1 * time.Millisecond) }
