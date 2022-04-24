package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/refurbedsender/gate"
	"github.com/dmalykh/refurbedsender/queue"
	"github.com/dmalykh/refurbedsender/sender"
	"sync"
	"time"
)

// Sender implements sender.Sender
type Sender struct {
	queue   queue.Queue
	gate    gate.Gate
	runOnce sync.Once
	// Errors
	errManager ErrorManager
	// Throttling
	rps          uint64
	tokens       uint64
	throttleOnce sync.Once
	lock         sync.Mutex
}

var ErrAddingToQueue = errors.New(`error`)

// The NewSender returns configured Sender
// If you
func NewSender(q queue.Queue, g gate.Gate, rps uint64, skipErrors bool) sender.Sender {
	var s = &Sender{
		queue:      q,
		gate:       g,
		errManager: newErrManager(skipErrors),
		rps:        rps,
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

// Run service only runOnce.
func (s *Sender) Run(ctx context.Context) error {
	var err error
	s.runOnce.Do(func() {
		go s.throttle(ctx, s.rps, 1, time.Duration(float64(time.Second)/float64(s.rps)))
		err = s.proceed(ctx, s.gate.Send)
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
		// Using waitgroup for control errors channel. Channel shouldn't be closed while all messages sent
		wg.Wait()
		s.shutdown()
	}()

	return s.queue.Consume(ctx, func(m sender.Message) {
		wg.Add(1)
		defer wg.Done()
		if err := s.throttlingMiddleware(ctx, f, m); err != nil {
			s.errManager.AddError(ctx, m, err)
		}
	})
}

func (s *Sender) shutdown() {
	s.errManager.Done()
}
