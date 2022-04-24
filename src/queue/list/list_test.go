package list

import (
	"context"
	"github.com/dmalykh/refurbedsender/sender"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestListQueue(t *testing.T) {
	var q = NewListQueue()
	ctx, cancel := context.WithCancel(context.TODO())
	var checks = 100

	var messages = make(map[uuid.UUID]struct{})
	for i := 0; i < checks*2; i++ {
		var msg = sender.NewMessage(``)
		messages[msg.GetID()] = struct{}{}
		assert.NoError(t, q.Add(ctx, msg))
	}

	var consumed int
	var lock sync.Mutex
	assert.NoError(t, q.Consume(ctx, func(message sender.Message) {
		lock.Lock()
		defer lock.Unlock()
		consumed++
		if consumed == checks {
			cancel()
		}
		_, ok := messages[message.GetID()]
		assert.True(t, ok)
	}))
	assert.Equal(t, checks, consumed)
}

func TestListQueueUnknownValue(t *testing.T) {
	var q = NewListQueue()
	q.list.PushBack(`strange value`)
	assert.True(t, assert.ErrorIs(t, q.Consume(context.TODO(), func(message sender.Message) {

	}), ErrUnknownValue))
}
