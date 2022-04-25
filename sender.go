package refurbedsender

import (
	"github.com/dmalykh/refurbedsender/gate"
	"github.com/dmalykh/refurbedsender/queue"
	"github.com/dmalykh/refurbedsender/sender"
	"github.com/dmalykh/refurbedsender/service"
)

func NewSender(q queue.Queue, g gate.Gate, skipErrors bool) sender.Sender {
	return service.NewSender(q, g, skipErrors)
}
