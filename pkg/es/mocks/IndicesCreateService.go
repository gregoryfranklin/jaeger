// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import es "github.com/jaegertracing/jaeger/pkg/es"
import mock "github.com/stretchr/testify/mock"

// IndicesCreateService is an autogenerated mock type for the IndicesCreateService type
type IndicesCreateService struct {
	mock.Mock
}

// Body provides a mock function with given fields: mapping
func (_m *IndicesCreateService) Body(mapping string) es.IndicesCreateService {
	ret := _m.Called(mapping)

	var r0 es.IndicesCreateService
	if rf, ok := ret.Get(0).(func(string) es.IndicesCreateService); ok {
		r0 = rf(mapping)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(es.IndicesCreateService)
		}
	}

	return r0
}

// Do provides a mock function with given fields: ctx
func (_m *IndicesCreateService) Do(ctx context.Context) (*es.IndicesCreateResult, error) {
	ret := _m.Called(ctx)

	var r0 *es.IndicesCreateResult
	if rf, ok := ret.Get(0).(func(context.Context) *es.IndicesCreateResult); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*es.IndicesCreateResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
