package mocks

import (
	"context"
	"ququiz/lintang/quiz-query-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockAuthGrpcClient struct {
	mock.Mock
}

func (_m *MockAuthGrpcClient) GetUsersByIds(ctx context.Context, userIDs []string) ([]domain.User, error) {
	ret := _m.Called(ctx, userIDs)
	if len(ret) == 0 {
		panic("Function GetUsersByIds has no return value")
	}

	var r0 []domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) []domain.User); ok {
		r0 = rf(ctx, userIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, userIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockAuthGrpcClient) GetUserByID(ctx context.Context, userID string) (domain.User, error) {
	ret := _m.Called(ctx, userID)
	if len(ret) == 0 {
		panic("Function GetUserByID has no return value")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

func NewAuthClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuthGrpcClient {
	m := &MockAuthGrpcClient{}
	m.Mock.Test(t)
	t.Cleanup(func() {
		m.AssertExpectations(t)
	})
	return m
}
