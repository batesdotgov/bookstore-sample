// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	purchases "github.com/diegoholiveira/bookstore-sample/purchases"
)

// PurchasePersister is an autogenerated mock type for the PurchasePersister type
type PurchasePersister struct {
	mock.Mock
}

// Persist provides a mock function with given fields: _a0, _a1
func (_m *PurchasePersister) Persist(_a0 context.Context, _a1 purchases.Purchase) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, purchases.Purchase) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
