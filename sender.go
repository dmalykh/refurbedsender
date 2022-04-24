package sender

import (
	"github.com/dmalykh/refurbedsender/gate"
	"github.com/dmalykh/refurbedsender/queue"
	"github.com/dmalykh/refurbedsender/sender"
	"github.com/dmalykh/refurbedsender/service"
)

func NewSender(q queue.Queue, g gate.Gate, rps uint64, skipErrors bool) sender.Sender {
	return service.NewSender(q, g, rps, skipErrors)
}
