package service

import (
	"context"
	"errors"
	"github.com/dmalykh/refurbedsender/sender"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrManager_Errors(t *testing.T) {
	type fields struct {
		skipErrors bool
	}
	tests := []struct {
		name   string
		fields fields
		panic  bool
		want   chan *sender.Error
	}{
		{
			`Panic when calling Errors() with skipErrors=true`,
			fields{
				skipErrors: true,
			},
			true,
			nil,
		},
		{
			`Channel returns`,
			fields{
				skipErrors: false,
			},
			false,
			make(chan *sender.Error),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newErrManager(tt.fields.skipErrors)
			if tt.panic {
				assert.Panics(t, func() {
					e.Errors()
				})

				return
			}
			assert.ObjectsAreEqual(tt.want, e.Errors())
			assert.IsType(t, tt.want, e.Errors())
		})
	}
}

func TestErrManager_AddError(t *testing.T) {
	type fields struct {
		skipErrors bool
	}
	type args struct {
		ctx context.Context //nolint:containedctx
		m   sender.Message
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			`Error should be skipped, Errors() panics`,
			fields{
				skipErrors: true,
			},
			args{
				ctx: context.TODO(),
				m:   sender.NewMessage(``),
				err: errors.New(`any`),
			},
		},
		{
			`Error added to channel`,
			fields{
				skipErrors: false,
			},
			args{
				ctx: context.TODO(),
				m:   sender.NewMessage(``),
				err: errors.New(`any`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errManager{
				skipErrors: tt.fields.skipErrors,
			}

			// Panic when calls Errors()
			if tt.fields.skipErrors {
				assert.Panics(t, func() {
					e.Errors()
				})
			}

			// Got error
			if !e.skipErrors {
				e.Err = make(chan *sender.Error)
				defer close(e.Err)

				go func() {
					got := <-e.Errors()
					assert.True(t, errors.Is(got.GetError(), tt.args.err))
					assert.Equal(t, tt.args.m.GetID(), got.GetMessageID())
				}()
			}

			e.AddError(tt.args.ctx, tt.args.m, tt.args.err)
		})
	}
}
