package service

import (
	"context"
	"github.com/dmalykh/refurbedsender/sender"
)

// ErrorManager interface used for working with errors in Sender
type ErrorManager interface {
	Errors() chan *sender.Error
	AddError(ctx context.Context, m sender.Message, err error)
	Done()
}

// errManager implements ErrorManager
type errManager struct {
	Err        chan *sender.Error
	skipErrors bool
}

// Constructor
func newErrManager(skipErrors bool) ErrorManager {
	var e = &errManager{
		skipErrors: skipErrors,
	}

	if !e.skipErrors {
		e.Err = make(chan *sender.Error)
	}

	return e
}

// Errors method returns channel with sending errors.
// Be afraid, when skipErrors equals true, this method returns nil.
func (e *errManager) Errors() chan *sender.Error {
	if e.skipErrors {
		panic(`errors shouldn't be used with skipErrors true`)
	}

	return e.Err
}

// AddError adds error message's id to errors channel
func (e *errManager) AddError(ctx context.Context, m sender.Message, err error) {
	if !e.skipErrors {
		e.Err <- &sender.Error{Err: err, MessageID: m.GetID()}
	}
}

// Done closes the errors channel
func (e *errManager) Done() {
	close(e.Err)
}
