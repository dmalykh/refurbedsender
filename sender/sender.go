package sender

import (
	"context"
)

type ProceedFunc func(ctx context.Context, message Message) error

type Sender interface {
	// The Run method starts queues for sending messages
	Run(ctx context.Context, opts ...Option) error
	// The Send method adds message to sending queue. Returns error if queue is not available
	Send(ctx context.Context, message Message) error
	// Errors method returns channel with errors of messages that was sent
	Errors() chan *Error
}
