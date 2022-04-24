package sender

import (
	"fmt"
	"github.com/google/uuid"
)

type Error struct {
	Err       error
	MessageID uuid.UUID
}

func (e *Error) Error() string {
	return fmt.Errorf(`(%s) %w`, e.GetMessageID().String(), e.Err).Error()
}

func (e *Error) GetMessageID() uuid.UUID {
	return e.MessageID
}

func (e *Error) GetError() error {
	return e.Err
}
