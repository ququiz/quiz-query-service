package mocks

import (
	"context"
	"ququiz/lintang/quiz-query-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockCacheRepo struct {
	mock.Mock
}

func (_m *MockCacheRepo) GetCachedQuestion(ctx context.Context, quizID string) ([]domain.Question, error) {

	ret := _m.Called(ctx, quizID)
	if len(ret) == 0 {
		panic("Function GetCachedQuestion has no return value")
	}

	var r0 []domain.Question
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) []domain.Question); ok {
		r0 = rf(ctx, quizID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Question)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, quizID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockCacheRepo) SetCachedQuestion(ctx context.Context, quizID string, qs []domain.Question) error {
	ret := _m.Called(ctx, quizID, qs)

	if len(ret) == 0 {
		panic("Function SetCachedQuestion has no return value")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []domain.Question) error); ok {
		r0 = rf(ctx, quizID, qs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}



func (_m *MockCacheRepo) DeleteCacheForSpecificQuiz(ctx context.Context, quizID string) error {
	ret := _m.Called(ctx, quizID)
	if len(ret) == 0 {
		panic("Function DeleteCacheForSpecificQuiz has no return value")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, quizID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func NewRedisCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCacheRepo {
	m := &MockCacheRepo{}
	m.Mock.Test(t)
	t.Cleanup(func() {
		m.AssertExpectations(t)
	})
	return m

}
