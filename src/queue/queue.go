package queue

import (
	"context"
	"github.com/dmalykh/refurbedsender/sender"
)

// Use queue as separate part. Because we may consider to use Kafka, MQ

type Queue interface {
	Add(ctx context.Context, message sender.Message) error
	Consume(ctx context.Context, f func(message sender.Message)) error
}
