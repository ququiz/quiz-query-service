package mocks

import (
	"context"
	"ququiz/lintang/quiz-query-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockQuizRepository struct {
	mock.Mock
}

func (_m *MockQuizRepository) InsertQuizData(ctx context.Context, quizReqs []domain.BaseQuiz) error {
	ret := _m.Called(ctx, quizReqs)

	if len(ret) == 0 {
		panic("Function InsertQuizData has no return value")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []domain.BaseQuiz) error); ok {
		r0 = rf(ctx, quizReqs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *MockQuizRepository) GetAll(ctx context.Context, limit uint64, offset uint64) ([]domain.BaseQuiz, error) {
	ret := _m.Called(ctx, limit, offset)
	if len(ret) == 0 {
		panic("Function GetAll has no return value")
	}

	var r0 []domain.BaseQuiz
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) []domain.BaseQuiz); ok {
		r0 = rf(ctx, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BaseQuiz)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, uint64) error); ok {
		r1 = rf(ctx, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockQuizRepository) Get(ctx context.Context, quizID string) (domain.BaseQuiz, error) {
	ret := _m.Called(ctx, quizID)

	if len(ret) == 0 {
		panic("Function Get has no return value")
	}

	var r0 domain.BaseQuiz
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.BaseQuiz); ok {
		r0 = rf(ctx, quizID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.BaseQuiz)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, quizID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockQuizRepository) IsUserQuizParticipant(ctx context.Context, quizID string, userID string) ([]domain.BaseQuizIsParticipant, error) {

	ret := _m.Called(ctx, quizID, userID)
	if len(ret) == 0 {
		panic("Function IsUserQuizParticipant has no return value")
	}

	var r0 []domain.BaseQuizIsParticipant
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []domain.BaseQuizIsParticipant); ok {
		r0 = rf(ctx, quizID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BaseQuizIsParticipant)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, quizID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockQuizRepository) GetAllQuizByCreatorID(ctx context.Context, creatorID string, limit uint64, offset uint64) ([]domain.BaseQuiz, error) {
	ret := _m.Called(ctx, creatorID, limit, offset)
	if len(ret) == 0 {
		panic("Function GetAllQuizByCreatorID has no return value")
	}

	var r0 []domain.BaseQuiz
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, uint64) []domain.BaseQuiz); ok {
		r0 = rf(ctx, creatorID, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BaseQuiz)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, uint64, uint64) error); ok {
		r1 = rf(ctx, creatorID, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockQuizRepository) GetQuizHistory(ctx context.Context, participantID string, limit uint64, offset uint64) ([]domain.BaseQuizIsParticipant, error) {
	ret := _m.Called(ctx, participantID, limit, offset)

	if len(ret) == 0 {
		panic("Function GetQuizHistory has no return value")
	}

	var r0 []domain.BaseQuizIsParticipant
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, uint64) []domain.BaseQuizIsParticipant); ok {
		r0 = rf(ctx, participantID, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BaseQuizIsParticipant)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, uint64, uint64) error); ok {
		r1 = rf(ctx, participantID, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func NewQuizRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockQuizRepository {
	m := &MockQuizRepository{}
	m.Mock.Test(t)
	t.Cleanup(func() {
		m.AssertExpectations(t)
	})
	return m
}
