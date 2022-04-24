package sender

import (
	"context"
)

type ProceedFunc func(ctx context.Context, message Message) error

type Sender interface {
	Run(ctx context.Context) error
	Send(ctx context.Context, message Message) error
	Errors() chan *Error
}
