package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/refurbedsender/gate"
	"github.com/dmalykh/refurbedsender/queue"
	"github.com/dmalykh/refurbedsender/sender"
	"sync"
)

// Sender implements sender.Sender
type Sender struct {
	queue   queue.Queue
	gate    gate.Gate
	runOnce sync.Once
	// Errors
	errManager ErrorManager
}

var ErrAddingToQueue = errors.New(`error`)

// The NewSender returns configured Sender
func NewSender(q queue.Queue, g gate.Gate, skipErrors bool) sender.Sender {
	var s = &Sender{
		queue:      q,
		gate:       g,
		errManager: newErrManager(skipErrors),
	}

	return s
}

// Send method used for adding message to sending queue, returns error if adding in queue is not available
func (s *Sender) Send(ctx context.Context, message sender.Message) error {
	if err := s.queue.Add(ctx, message); err != nil {
		return fmt.Errorf(`%w: %s`, ErrAddingToQueue, err.Error())
	}

	return nil
}

// Run service only runOnce. Use opts for settings middlewares
func (s *Sender) Run(ctx context.Context, opts ...sender.Option) error {
	var err error
	s.runOnce.Do(func() {
		var handler = s.gate.Send
		for _, opt := range opts {
			handler = opt(handler)
		}
		err = s.proceed(ctx, handler)
	})

	return err
}

func (s *Sender) Errors() chan *sender.Error {
	return s.errManager.Errors()
}

//
func (s *Sender) proceed(ctx context.Context, f sender.ProceedFunc) error {
	var wg = new(sync.WaitGroup)
	defer func() {
		// Using waitgroup for control shutdown
		wg.Wait()
		s.shutdown()
	}()

	return s.queue.Consume(ctx, func(m sender.Message) {
		go func() {
			wg.Add(1)
			defer wg.Done()
			if err := f(ctx, m); err != nil {
				s.errManager.AddError(ctx, m, err)
			}
		}()
	})
}

func (s *Sender) shutdown() {
	s.errManager.Done()
}
