package mocks

import (
	"context"
	"ququiz/lintang/quiz-query-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockQuestionRepository struct {
	mock.Mock
}

func (_m *MockQuestionRepository) Get(ctx context.Context, questionID string) (domain.Question, error) {
	ret := _m.Called(ctx, questionID)
	if len(ret) == 0 {
		panic("Function Get has no return value")
	}

	var r0 domain.Question
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Question); ok {
		r0 = rf(ctx, questionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Question)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, questionID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

func (_m *MockQuestionRepository) GetAllByQuiz(ctx context.Context, quizID string) ([]domain.BaseQuizWithQuestionAggregate, error) {
	ret := _m.Called(ctx, quizID)
	if len(ret) == 0 {
		panic("Function GetAllByQuiz has no return value")
	}

	var r0 []domain.BaseQuizWithQuestionAggregate
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) []domain.BaseQuizWithQuestionAggregate); ok {
		r0 = rf(ctx, quizID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BaseQuizWithQuestionAggregate)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, quizID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

func (_m *MockQuestionRepository) GetUserAnswerInAQuiz(ctx context.Context, quizID string, userID string) ([]domain.QuestionUserAnswer, error) {
	ret := _m.Called(ctx, quizID, userID)
	if len(ret) == 0 {
		panic("Function GetUserAnswerInAQuiz has no return value")
	}

	var r0 []domain.QuestionUserAnswer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []domain.QuestionUserAnswer); ok {
		r0 = rf(ctx, quizID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.QuestionUserAnswer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, quizID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockQuestionRepository) IsUserAlreadyAnswerThisQuizID(ctx context.Context, quizID string,
	questionID string,
	userID string) (bool, error) {

	ret := _m.Called(ctx, quizID, questionID, userID)
	if len(ret) == 0 {
		panic("Function IsUserAlreadyAnswerThisQuizID has no return value")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) bool); ok {
		r0 = rf(ctx, quizID, questionID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(bool)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, quizID, questionID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockQuestionRepository) IsUserAnswerCorrect(ctx context.Context, quizID string, questionID string,
	userChoiceID string, userEssayAnswer string) (bool, domain.CorrectAnswer, error) {
	ret := _m.Called(ctx, quizID, questionID, userChoiceID, userEssayAnswer)
	if len(ret) == 0 {
		panic("Function IsUserAnswerCorrect has no return value")
	}

	var r0 bool
	var r1 domain.CorrectAnswer
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) bool); ok {
		r0 = rf(ctx, quizID, questionID, userChoiceID, userEssayAnswer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(bool)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string) domain.CorrectAnswer); ok {
		r1 = rf(ctx, quizID, questionID, userChoiceID, userEssayAnswer)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(domain.CorrectAnswer)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, string, string) error); ok {
		r2 = rf(ctx, quizID, questionID, userChoiceID, userEssayAnswer)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2

}

func (_m *MockQuestionRepository) GetQuestionByIDAndQuizID(ctx context.Context, quizID string, questionID string) (domain.BaseQuizWithOneQuestionAggregate, error) {
	ret := _m.Called(ctx, quizID, questionID)

	if len(ret) == 0 {
		panic("Function GetQuestionByIDAndQuizID has no return value")
	}

	var r0 domain.BaseQuizWithOneQuestionAggregate
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) domain.BaseQuizWithOneQuestionAggregate); ok {
		r0 = rf(ctx, quizID, questionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.BaseQuizWithOneQuestionAggregate)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, quizID, questionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func NewQuestionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockQuestionRepository {
	m := &MockQuestionRepository{}
	m.Mock.Test(t)
	t.Cleanup(func() {
		m.AssertExpectations(t)
	})
	return m
}
