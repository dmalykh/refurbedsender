package middleware

import (
	"context"
	"github.com/dmalykh/refurbedsender/sender"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestThrottlingMiddleware(t *testing.T) {
	var call = 1000
	var wantCalls = 1000
	var called int32
	var f sender.ProceedFunc = func(ctx context.Context, message sender.Message) error {
		atomic.AddInt32(&called, 1)

		return nil
	}
	var throttle = ThrottlingMiddleware(f, 2, 1, time.Millisecond)
	var ctx = context.TODO()
	for i := 0; i < call; i++ {
		assert.NoError(t, throttle(ctx, sender.NewMessage(``)))
	}
	assert.Equal(t, wantCalls, int(called))
}

func TestThrottlingMiddlewareStopContext(t *testing.T) {
	var call = 100
	var wantCalls = 6
	var called int32
	var f sender.ProceedFunc = func(ctx context.Context, message sender.Message) error {
		atomic.AddInt32(&called, 1)

		return nil
	}
	var throttle = ThrottlingMiddleware(f, 4, 2, time.Second)
	//goland:noinspection GoVetLostCancel
	ctx, _ := context.WithTimeout(context.TODO(), 1500*time.Millisecond) // 1.5s
	for i := 0; i < call; i++ {
		go func() {
			err := throttle(ctx, sender.NewMessage(``))
			if err != nil {
				assert.ErrorIs(t, err, ctx.Err())
			}
		}()
	}
	<-ctx.Done()
	assert.Equal(t, wantCalls, int(atomic.LoadInt32(&called)))
}
