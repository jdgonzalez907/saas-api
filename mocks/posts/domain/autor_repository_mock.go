package mock_domain

import (
	"context"

	"github.com/jdgonzalez907/saas-api/internal/posts/domain"
	mock "github.com/stretchr/testify/mock"
)

func NewMockAutorRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAutorRepository {
	mock := &MockAutorRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

type MockAutorRepository struct {
	mock.Mock
}

type MockAutorRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAutorRepository) EXPECT() *MockAutorRepository_Expecter {
	return &MockAutorRepository_Expecter{mock: &_m.Mock}
}

func (_mock *MockAutorRepository) FindByID(ctx context.Context, id int64) (*domain.Autor, error) {
	ret := _mock.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 *domain.Autor
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, int64) (*domain.Autor, error)); ok {
		return returnFunc(ctx, id)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, int64) *domain.Autor); ok {
		r0 = returnFunc(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Autor)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = returnFunc(ctx, id)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

type MockAutorRepository_FindByID_Call struct {
	*mock.Call
}

func (_e *MockAutorRepository_Expecter) FindByID(ctx any, id any) *MockAutorRepository_FindByID_Call {
	return &MockAutorRepository_FindByID_Call{Call: _e.mock.On("FindByID", ctx, id)}
}

func (_c *MockAutorRepository_FindByID_Call) Run(run func(ctx context.Context, id int64)) *MockAutorRepository_FindByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 int64
		if args[1] != nil {
			arg1 = args[1].(int64)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockAutorRepository_FindByID_Call) Return(autor *domain.Autor, err error) *MockAutorRepository_FindByID_Call {
	_c.Call.Return(autor, err)
	return _c
}

func (_c *MockAutorRepository_FindByID_Call) RunAndReturn(run func(ctx context.Context, id int64) (*domain.Autor, error)) *MockAutorRepository_FindByID_Call {
	_c.Call.Return(run)
	return _c
}
