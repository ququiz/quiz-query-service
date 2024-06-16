package mocks

import (
	"context"
	"ququiz/lintang/quiz-query-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockQuizCommandProducer struct {
	mock.Mock
}

func (_m *MockQuizCommandProducer) SendCorrectAnswerToQuizCommandService(ctx context.Context, userAnswerMsg domain.UserAnswerMQ) error {
	ret := _m.Called(ctx, userAnswerMsg)
	if len(ret) == 0 {
		panic("Function SendCorrectAnswerToQuizCommandService has no return value")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserAnswerMQ) error); ok {
		r0 = rf(ctx, userAnswerMsg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func NewQuizCommandServiceProducerMQ(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockQuizCommandProducer {
	m := &MockQuizCommandProducer{}
	m.Mock.Test(t)

	t.Cleanup(func() {
		m.AssertExpectations(t)
	})
	return m
}

