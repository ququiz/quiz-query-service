package mocks

import (
	"context"
	"ququiz/lintang/quiz-query-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockScoringProducer struct {
	mock.Mock
}

func (_m *MockScoringProducer) SendCorrectAnswer(ctx context.Context, correctAnswerMsg domain.CorrectAnswer) error {
	ret := _m.Called(ctx, correctAnswerMsg)

	if len(ret) == 0 {
		panic("Function SendCorrectAnswer has no return value")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.CorrectAnswer) error); ok {
		r0 = rf(ctx, correctAnswerMsg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func NewScoringServiceProducerMQ(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockScoringProducer {
	m := &MockScoringProducer{}
	m.Mock.Test(t)

	t.Cleanup(func() {
		m.AssertExpectations(t)
	})
	return m
}
