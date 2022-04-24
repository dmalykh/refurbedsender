package service

import (
	"context"
	"errors"
	"github.com/dmalykh/refurbedsender/mocks"
	"github.com/dmalykh/refurbedsender/sender"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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
