// Code generated by mockery v2.12.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sender "github.com/dmalykh/refurbedsender/sender"

	testing "testing"
)

// Queue is an autogenerated mock type for the Queue type
type Queue struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, message
func (_m *Queue) Add(ctx context.Context, message sender.Message) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, sender.Message) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Consume provides a mock function with given fields: ctx, f
func (_m *Queue) Consume(ctx context.Context, f func(sender.Message)) error {
	ret := _m.Called(ctx, f)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(sender.Message)) error); ok {
		r0 = rf(ctx, f)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewQueue creates a new instance of Queue. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewQueue(t testing.TB) *Queue {
	mock := &Queue{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
