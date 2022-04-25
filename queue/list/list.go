package list

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/refurbedsender/sender"
	"reflect"
)

var ErrUnknownValue = errors.New(`received unknown value in lists element`)

func NewListQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

// Queue is a simple queue based on linked list.
type Queue struct {
	list *list.List
}

func (l *Queue) Add(ctx context.Context, message sender.Message) error {
	l.list.PushBack(message)

	return nil
}

func (l *Queue) Consume(ctx context.Context, f func(message sender.Message)) error {
	for {
		select {
		case <-ctx.Done():
			l.list.Init()

			return nil
		default:
			if element := l.list.Front(); element != nil {
				message, ok := element.Value.(sender.Message)
				if !ok {
					return fmt.Errorf(`%w: %s`, ErrUnknownValue, reflect.TypeOf(element.Value))
				}
				f(message)
				l.list.Remove(element)
			}
		}
	}
}
