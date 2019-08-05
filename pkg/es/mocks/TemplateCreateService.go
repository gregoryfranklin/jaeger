// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import elastic "github.com/olivere/elastic"
import es "github.com/jaegertracing/jaeger/pkg/es"
import mock "github.com/stretchr/testify/mock"

// TemplateCreateService is an autogenerated mock type for the TemplateCreateService type
type TemplateCreateService struct {
	mock.Mock
}

// Body provides a mock function with given fields: mapping
func (_m *TemplateCreateService) Body(mapping string) es.TemplateCreateService {
	ret := _m.Called(mapping)

	var r0 es.TemplateCreateService
	if rf, ok := ret.Get(0).(func(string) es.TemplateCreateService); ok {
		r0 = rf(mapping)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(es.TemplateCreateService)
		}
	}

	return r0
}

// Do provides a mock function with given fields: ctx
func (_m *TemplateCreateService) Do(ctx context.Context) (*elastic.IndicesPutTemplateResponse, error) {
	ret := _m.Called(ctx)

	var r0 *elastic.IndicesPutTemplateResponse
	if rf, ok := ret.Get(0).(func(context.Context) *elastic.IndicesPutTemplateResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*elastic.IndicesPutTemplateResponse)
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
