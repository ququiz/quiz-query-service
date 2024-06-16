package service

import (
	"context"
	"errors"
	"fmt"
	"ququiz/lintang/quiz-query-service/biz/domain"
	"ququiz/lintang/quiz-query-service/biz/service/mocks"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllByQuiz(t *testing.T) {
	mockQuizRepo := new(mocks.MockQuizRepository)
	mockQuestionRepo := new(mocks.MockQuestionRepository)
	mockCacheRepo := new(mocks.MockCacheRepo)
	mockScoringProducer := new(mocks.MockScoringProducer)
	mockQuizCommandProducer := new(mocks.MockQuizCommandProducer)

	var mockQuizID = "666d8faaed25031b0d947430"
	var mockUserID = "666d8faaed25031b0d91234"

	var mockQuiz domain.BaseQuiz
	err := faker.FakeData(&mockQuiz)
	assert.NoError(t, err)
	var mockQuestion = make([]domain.Question, 0)
	for i := 0; i < 3; i++ {
		var question domain.Question
		err = faker.FakeData(&question)

		mockQuestion = append(mockQuestion, question)
	}

	var mockParticipants = make([]domain.BaseQuizIsParticipant, 0)
	var mockParticipant domain.BaseQuizIsParticipant
	err = faker.FakeData(&mockParticipant)
	assert.NoError(t, err)
	mockParticipants = append(mockParticipants, mockParticipant)

	var mockQuestionFromDB []domain.BaseQuizWithQuestionAggregate = make([]domain.BaseQuizWithQuestionAggregate, 0)

	for i := 0; i < 2; i++ {
		var mockQuestionDB domain.BaseQuizWithQuestionAggregate
		err = faker.FakeData(&mockQuestionDB)
		assert.NoError(t, err)
		mockQuestionFromDB = append(mockQuestionFromDB, mockQuestionDB)
	}

	t.Run("success get question from cache", func(t *testing.T) {
		mockQuiz.StartTime = time.Now().Add(-time.Hour * 2)
		mockQuiz.EndTime = time.Now().Add(time.Hour * 2)

		mockQuizRepo.On("IsUserQuizParticipant", mock.Anything, mockQuizID, mockUserID).Return(mockParticipants, nil).Once()
		mockQuizRepo.On("Get", mock.Anything, mockQuizID).Return(mockQuiz, nil).Once()

		mockCacheRepo.On("GetCachedQuestion", mock.Anything, mockQuizID).Return(mockQuestion, nil).Once()

		service := NewQuestionService(mockQuestionRepo, mockCacheRepo, mockQuizRepo, mockScoringProducer, mockQuizCommandProducer)
		questions, err := service.GetAllByQuiz(context.TODO(), mockQuizID, mockUserID)

		assert.NoError(t, err)
		assert.NotEmpty(t, questions)
		assert.Equal(t, mockQuestion, questions)
		assert.Len(t, questions, len(mockQuestion))

		mockQuizRepo.AssertExpectations(t)
		mockQuizRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)

	})

	t.Run("success get question from db", func(t *testing.T) {
		mockQuiz.StartTime = time.Now().Add(-time.Hour * 2)
		mockQuiz.EndTime = time.Now().Add(time.Hour * 2)

		mockQuizRepo.On("IsUserQuizParticipant", mock.Anything, mockQuizID, mockUserID).Return(mockParticipants, nil).Once()
		mockQuizRepo.On("Get", mock.Anything, mockQuizID).Return(mockQuiz, nil).Once()

		mockCacheRepo.On("GetCachedQuestion", mock.Anything, mockQuizID).Return(nil, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, fmt.Sprintf(`quiz not found`))).Once()

		mockQuestionRepo.On("GetAllByQuiz", mock.Anything, mockQuizID).Return(mockQuestionFromDB, nil).Once()

		var cachedQuestions []domain.Question = make([]domain.Question, 0)
		for _, quiz := range mockQuestionFromDB {
			cachedQuestions = append(cachedQuestions, quiz.Questions...)
		}

		mockCacheRepo.On("SetCachedQuestion", mock.Anything, mockQuizID, cachedQuestions).Return(nil).Once()

		service := NewQuestionService(mockQuestionRepo, mockCacheRepo, mockQuizRepo, mockScoringProducer, mockQuizCommandProducer)
		questions, err := service.GetAllByQuiz(context.TODO(), mockQuizID, mockUserID)

		assert.NoError(t, err)
		assert.NotEmpty(t, questions)
		assert.Equal(t, questions, cachedQuestions)
		assert.Len(t, questions, len(cachedQuestions))

		mockQuizRepo.AssertExpectations(t)
		mockQuizRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
		mockQuestionRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})

}

func TestGetUserAnswer(t *testing.T) {
	mockQuizRepo := new(mocks.MockQuizRepository)
	mockQuestionRepo := new(mocks.MockQuestionRepository)
	mockCacheRepo := new(mocks.MockCacheRepo)
	mockScoringProducer := new(mocks.MockScoringProducer)
	mockQuizCommandProducer := new(mocks.MockQuizCommandProducer)

	var mockQuizID = "666d8faaed25031b0d947430"
	var mockUserID = "666d8faaed25031b0d91234"

	var mockQuiz domain.BaseQuiz
	err := faker.FakeData(&mockQuiz)
	assert.NoError(t, err)
	var mockQuestion = make([]domain.Question, 0)
	for i := 0; i < 3; i++ {
		var question domain.Question
		err = faker.FakeData(&question)

		mockQuestion = append(mockQuestion, question)
	}

	var mockParticipants = make([]domain.BaseQuizIsParticipant, 0)
	var mockParticipant domain.BaseQuizIsParticipant
	err = faker.FakeData(&mockParticipant)
	assert.NoError(t, err)
	mockParticipants = append(mockParticipants, mockParticipant)

	var mockQuestionFromDB []domain.BaseQuizWithQuestionAggregate = make([]domain.BaseQuizWithQuestionAggregate, 0)

	for i := 0; i < 2; i++ {
		var mockQuestionDB domain.BaseQuizWithQuestionAggregate
		err = faker.FakeData(&mockQuestionDB)
		assert.NoError(t, err)
		mockQuestionFromDB = append(mockQuestionFromDB, mockQuestionDB)
	}

	var userAnswer []domain.QuestionUserAnswer = make([]domain.QuestionUserAnswer, 0)
	for i := 0; i < 3; i++ {
		var answer domain.QuestionUserAnswer
		err = faker.FakeData(&answer)
		assert.NoError(t, err)
		userAnswer = append(userAnswer, answer)
	}

	t.Run("success get user answer ", func(t *testing.T) {
		mockQuiz.StartTime = time.Now().Add(-time.Hour * 2)
		mockQuiz.EndTime = time.Now().Add(time.Hour * 2)

		mockQuizRepo.On("IsUserQuizParticipant", mock.Anything, mockQuizID, mockUserID).Return(mockParticipants, nil).Once()

		mockQuestionRepo.On("GetUserAnswerInAQuiz", mock.Anything, mockQuizID, mockUserID).Return(userAnswer, nil).Once()

		service := NewQuestionService(mockQuestionRepo, mockCacheRepo, mockQuizRepo, mockScoringProducer, mockQuizCommandProducer)
		userAnswersRes, err := service.GetUserAnswers(context.TODO(), mockQuizID, mockUserID)

		assert.NoError(t, err)
		assert.NotEmpty(t, userAnswersRes)
		assert.Equal(t, userAnswersRes, userAnswer)
		assert.Len(t, userAnswersRes, len(userAnswer))

		mockQuizRepo.AssertExpectations(t)
		mockQuestionRepo.AssertExpectations(t)
	})
}

func TestUserAnswerAQuestion(t *testing.T) {
	mockQuizRepo := new(mocks.MockQuizRepository)
	mockQuestionRepo := new(mocks.MockQuestionRepository)
	mockCacheRepo := new(mocks.MockCacheRepo)
	mockScoringProducer := new(mocks.MockScoringProducer)
	mockQuizCommandProducer := new(mocks.MockQuizCommandProducer)

	var mockQuizID = "666d8faaed25031b0d947430"
	var mockUserID = "666d8faaed25031b0d91234"
	var mockQuestionID = "666d8faaed25031b0d91234"
	var mockUserChoiceID = "666d8faaed25031b0d91234"
	var mockUserEssayAnswer = "666d8faaed25031b0d91234"
	var mockUsername = "lintang"

	var mockQuiz domain.BaseQuiz
	err := faker.FakeData(&mockQuiz)
	assert.NoError(t, err)
	var mockQuestion = make([]domain.Question, 0)
	for i := 0; i < 3; i++ {
		var question domain.Question
		err = faker.FakeData(&question)

		mockQuestion = append(mockQuestion, question)
	}

	var mockParticipants = make([]domain.BaseQuizIsParticipant, 0)
	var mockParticipant domain.BaseQuizIsParticipant
	err = faker.FakeData(&mockParticipant)
	assert.NoError(t, err)
	mockParticipants = append(mockParticipants, mockParticipant)

	var mockQuestionFromDB []domain.BaseQuizWithQuestionAggregate = make([]domain.BaseQuizWithQuestionAggregate, 0)

	for i := 0; i < 2; i++ {
		var mockQuestionDB domain.BaseQuizWithQuestionAggregate
		err = faker.FakeData(&mockQuestionDB)
		assert.NoError(t, err)
		mockQuestionFromDB = append(mockQuestionFromDB, mockQuestionDB)
	}

	var userAnswer []domain.QuestionUserAnswer = make([]domain.QuestionUserAnswer, 0)
	for i := 0; i < 3; i++ {
		var answer domain.QuestionUserAnswer
		err = faker.FakeData(&answer)
		assert.NoError(t, err)
		userAnswer = append(userAnswer, answer)
	}

	t.Run("success user answer a question and the answer is right ", func(t *testing.T) {
		mockQuiz.StartTime = time.Now().Add(-time.Hour * 2)
		mockQuiz.EndTime = time.Now().Add(time.Hour * 2)

		mockQuizRepo.On("IsUserQuizParticipant", mock.Anything, mockQuizID, mockUserID).Return(mockParticipants, nil).Once()

		mockQuestionRepo.On("IsUserAlreadyAnswerThisQuizID", mock.Anything, mockQuizID, mockQuestionID, mockUserID).Return(false, nil).Once()

		var mockCorrectAnswer domain.CorrectAnswer
		err = faker.FakeData(&mockCorrectAnswer)
		assert.NoError(t, err)
		mockCorrectAnswer.UserID = mockUserID
		mockCorrectAnswer.Username = mockUsername

		mockQuestionRepo.On("IsUserAnswerCorrect", mock.Anything, mockQuizID, mockQuestionID, mockUserChoiceID, mockUserEssayAnswer).Return(true, mockCorrectAnswer, nil).Once()

		mockScoringProducer.On("SendCorrectAnswer", mock.Anything, mockCorrectAnswer).Return(nil).Once()

		var userAnswerMQ domain.UserAnswerMQ = domain.UserAnswerMQ{
			QuizID:        mockQuizID,
			QuestionID:    mockQuestionID,
			ChoiceID:      mockUserChoiceID,
			ParticipantID: mockUserID,
			Answer:        mockUserEssayAnswer,
		}

		mockQuizCommandProducer.On("SendCorrectAnswerToQuizCommandService", mock.Anything, userAnswerMQ).Return(nil).Once()

		service := NewQuestionService(mockQuestionRepo, mockCacheRepo, mockQuizRepo, mockScoringProducer, mockQuizCommandProducer)
		isCorrect, err := service.UserAnswerAQuestion(context.TODO(), mockQuizID, mockUserID, mockQuestionID, mockUserChoiceID, mockUserEssayAnswer, mockUsername)

		assert.NoError(t, err)

		assert.Equal(t, true, isCorrect)

		mockQuizRepo.AssertExpectations(t)
		mockQuestionRepo.AssertExpectations(t)
		mockQuestionRepo.AssertExpectations(t)
		mockScoringProducer.AssertExpectations(t)
		mockQuizCommandProducer.AssertExpectations(t)
	})

}
