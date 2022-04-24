package service

import (
	"context"
	"errors"
	"github.com/dmalykh/refurbedsender/sender"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestSender_throttlingMiddlewareContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	var countCalls = 10

	// Make throttling func
	var check = func(ctx context.Context, m sender.Message) error {
		time.Sleep(1 * time.Second)

		return errors.New(`error`)
	}

	s := &Sender{}
	go s.throttle(ctx, uint64(countCalls), 1, 1*time.Minute)
	for i := 0; i < countCalls; i++ {
		go func() {
			err := s.throttlingMiddleware(ctx, check, sender.NewMessage(``))
			assert.NoError(t, err)
		}()
	}
	time.Sleep(100 * time.Millisecond)
	cancel()
}

func TestSender_throttlingMiddlewareSyncReceiving(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	var countCalls = 10

	// Make throttling func
	var check = func(ctx context.Context, m sender.Message) error {
		time.Sleep(100 * time.Millisecond)

		return errors.New(m.GetID().String())
	}

	// Run middlewares
	s := &Sender{}
	go s.throttle(ctx, uint64(countCalls), 1, 1*time.Minute)
	var wg = new(sync.WaitGroup)
	for i := 0; i < countCalls; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var message = sender.NewMessage(`not matter`)
			err := s.throttlingMiddleware(ctx, check, message)
			assert.Error(t, err)
			assert.Equal(t, err.Error(), message.GetID().String())
		}(i)
	}
	wg.Wait()
	cancel()
}

// The for loop in the throttlingMiddleware should call any ProceedFunc only once.
func TestSender_throttlingMiddlewareFuncCalledOnce(t *testing.T) {
	var ctx = context.TODO()
	var countCalls = 10

	var called int
	var check = func(ctx context.Context, m sender.Message) error {
		called++
		time.Sleep(100 * time.Millisecond)

		return nil
	}
	s := &Sender{}
	go s.throttle(ctx, 10, 1, 1*time.Millisecond)
	var wg = new(sync.WaitGroup)
	for i := 0; i < countCalls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(t, s.throttlingMiddleware(ctx, check, sender.NewMessage("test")))
		}()
	}
	wg.Wait()
	assert.Equal(t, countCalls, called)
}
