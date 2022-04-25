package service

import (
	"context"
	"errors"
	"github.com/dmalykh/refurbedsender/mocks"
	"github.com/dmalykh/refurbedsender/sender"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestSender_Send(t *testing.T) {
	tests := []struct {
		name      string
		testCalls int
		err       error
	}{
		{
			`No errors`,
			10,
			nil,
		},
		{
			`Error to add to queue`,
			10,
			errors.New(``),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.TODO()
			var queue = mocks.NewQueue(t)
			queue.On(`Add`, mock.Anything, mock.Anything).Return(func(ctx context.Context, m sender.Message) error {
				return tt.err
			})

			var s = &Sender{
				queue: queue,
			}
			for i := 0; i < tt.testCalls; i++ {
				if tt.err != nil {
					assert.Error(t, s.Send(ctx, sender.NewMessage(``)))
				} else {
					assert.NoError(t, s.Send(ctx, sender.NewMessage(``)))
				}
			}
			queue.AssertNumberOfCalls(t, `Add`, tt.testCalls)
		})
	}
}

// Run created once and middleware called
func TestSender_Run(t *testing.T) {
	var setOptCalled int
	var middlewareCalled int

	var middleware sender.Option = func(f sender.ProceedFunc) sender.ProceedFunc {
		setOptCalled++

		return func(ctx context.Context, message sender.Message) error {
			middlewareCalled++

			return nil
		}
	}

	var ctx = context.TODO()
	var q = mocks.NewQueue(t)
	q.On(`Consume`, mock.Anything, mock.Anything).Return(func(ctx context.Context, f func(message sender.Message)) error {
		f(sender.NewMessage(``))

		return nil
	})
	var s = NewSender(q, mocks.NewGate(t), true)
	assert.NoError(t, s.Run(ctx, middleware))
	assert.NoError(t, s.Run(ctx, middleware))
	assert.NoError(t, s.Run(ctx, middleware))
	assert.NoError(t, s.Run(ctx, middleware))

	assert.Equal(t, 1, setOptCalled)
	assert.Equal(t, 1, middlewareCalled)
}

func TestSender_proceedErrorAdded(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			`Error called`,
			errors.New(`ErrInFunc`),
		},
		{
			`No error`,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var message = sender.NewMessage(``)
			var middleware sender.Option = func(f sender.ProceedFunc) sender.ProceedFunc {
				return func(ctx context.Context, message sender.Message) error {
					time.Sleep(100 * time.Millisecond)

					return tt.err
				}
			}

			var ctx = context.TODO()
			var q = mocks.NewQueue(t)
			q.On(`Consume`, mock.Anything, mock.Anything).Return(func(ctx context.Context, f func(message sender.Message)) error {
				f(message)

				return nil
			})

			var s = NewSender(q, mocks.NewGate(t), false)
			go func() {
				assert.NoError(t, s.Run(ctx, middleware))
			}()

			for err := range s.Errors() {
				assert.ErrorIs(t, err.GetError(), tt.err)
				assert.Equal(t, message.GetID(), err.GetMessageID())
			}
		})
	}
}
