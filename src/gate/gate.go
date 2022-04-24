package gate

import (
	"context"
	"github.com/dmalykh/refurbedsender/sender"
)

type Gate interface {
	Send(ctx context.Context, messages sender.Message) error
}
