// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	mailer "github.com/tsel-ticketmaster/tm-notification/pkg/mailer"
)

// Mailer is an autogenerated mock type for the Mailer type
type Mailer struct {
	mock.Mock
}

// Send provides a mock function with given fields: ctx, messages
func (_m *Mailer) Send(ctx context.Context, messages ...mailer.Message) error {
	_va := make([]interface{}, len(messages))
	for _i := range messages {
		_va[_i] = messages[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...mailer.Message) error); ok {
		r0 = rf(ctx, messages...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
