package service_test

import (
	"context"
	"errors"
	"fmt"
	"ququiz/lintang/quiz-query-service/biz/domain"
	"ququiz/lintang/quiz-query-service/biz/service"
	"ququiz/lintang/quiz-query-service/biz/service/mocks"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetQuizs(t *testing.T) {
	mockQuizRepo := new(mocks.MockQuizRepository)
	mockAuthClient := new(mocks.MockAuthGrpcClient)
	mockLimit := 10
	mockOffset := 0

	mockListQuiz := make([]domain.BaseQuiz, 0)

	for i := 0; i < 10; i++ {
		var newMockFakerQuiz domain.BaseQuiz
		err := faker.FakeData(&newMockFakerQuiz)
		assert.NoError(t, err)
		mockListQuiz = append(mockListQuiz, newMockFakerQuiz)
	}

	mockListUser := make([]domain.User, 0)
	for i := 0; i < 10; i++ {
		var newMockFakerUser domain.User
		err := faker.FakeData(&newMockFakerUser)
		assert.NoError(t, err)
		mockListUser = append(mockListUser, newMockFakerUser)
	}

	// hasil quiz akhir

	t.Run("success get quiz", func(t *testing.T) {

		mockQuizRepo.On("GetAll", mock.Anything, uint64(mockLimit), uint64(mockOffset)).Return(mockListQuiz, nil).Once()

		quizService := service.NewQuizService(mockQuizRepo, mockAuthClient)

		var userIDs []string
		for i := 0; i < len(mockListQuiz); i++ {

			userIDs = append(userIDs, mockListQuiz[i].CreatorID)
		}

		// get user detail
		mockAuthClient.On("GetUsersByIds", mock.Anything, userIDs).Return(mockListUser, nil).Once()

		var userIDmap map[string]string = make(map[string]string)
		for i := 0; i < len(mockListUser); i++ {
			userIDmap[mockListUser[i].ID] = mockListUser[i].Username
		}

		for i := 0; i < len(mockListQuiz); i++ {
			mockListQuiz[i].CreatorName = userIDmap[mockListQuiz[i].CreatorID]
		}

		quizs, err := quizService.GetAll(context.TODO(), uint64(mockLimit), uint64(mockOffset))

		assert.NoError(t, err)
		assert.NotEmpty(t, quizs)
		assert.Equal(t, quizs, mockListQuiz)
		assert.Len(t, quizs, len(mockListQuiz))

		mockQuizRepo.AssertExpectations(t)
		mockAuthClient.AssertExpectations(t)
	})

	t.Run("not found when get quiz", func(t *testing.T) {

		mockQuizRepo.On("GetAll", mock.Anything, uint64(mockLimit), uint64(mockOffset)).Return([]domain.BaseQuiz{}, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, fmt.Sprintf(`quiz not found`))).Once()

		quizService := service.NewQuizService(mockQuizRepo, mockAuthClient)

		var userIDs []string
		for i := 0; i < len(mockListQuiz); i++ {

			userIDs = append(userIDs, mockListQuiz[i].CreatorID)
		}

		// get user detail
		mockAuthClient.On("GetUsersByIds", mock.Anything, userIDs).Return([]domain.User{}, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, fmt.Sprintf(`quiz not found`))).Once()

		var userIDmap map[string]string = make(map[string]string)
		for i := 0; i < len(mockListUser); i++ {
			userIDmap[mockListUser[i].ID] = mockListUser[i].Username
		}

		for i := 0; i < len(mockListQuiz); i++ {
			mockListQuiz[i].CreatorName = userIDmap[mockListQuiz[i].CreatorID]
		}

		quizs, err := quizService.GetAll(context.TODO(), uint64(mockLimit), uint64(mockOffset))

		assert.Error(t, err)
		assert.Equal(t, err, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, fmt.Sprintf(`quiz not found`)))
		assert.Empty(t, quizs)

		mockQuizRepo.AssertExpectations(t)
	})

	t.Run("error-failed when get quiz", func(t *testing.T) {

		mockQuizRepo.On("GetAll", mock.Anything, uint64(mockLimit), uint64(mockOffset)).Return([]domain.BaseQuiz{}, domain.WrapErrorf(errors.New(""), domain.ErrInternalServerError, fmt.Sprintf(`Unexpected Error`))).Once()

		quizService := service.NewQuizService(mockQuizRepo, mockAuthClient)

		var userIDs []string
		for i := 0; i < len(mockListQuiz); i++ {

			userIDs = append(userIDs, mockListQuiz[i].CreatorID)
		}

		// get user detail
		mockAuthClient.On("GetUsersByIds", mock.Anything, userIDs).Return([]domain.User{}, domain.ErrInternalServerError).Once()

		var userIDmap map[string]string = make(map[string]string)
		for i := 0; i < len(mockListUser); i++ {
			userIDmap[mockListUser[i].ID] = mockListUser[i].Username
		}

		for i := 0; i < len(mockListQuiz); i++ {
			mockListQuiz[i].CreatorName = userIDmap[mockListQuiz[i].CreatorID]
		}

		quizs, err := quizService.GetAll(context.TODO(), uint64(mockLimit), uint64(mockOffset))

		assert.Error(t, err)
		assert.Equal(t, err, domain.WrapErrorf(errors.New(""), domain.ErrInternalServerError, fmt.Sprintf(`Unexpected Error`)))
		assert.Empty(t, quizs)

		mockQuizRepo.AssertExpectations(t)
	})
}

func TestGetQuiz(t *testing.T) {
	mockQuizRepo := new(mocks.MockQuizRepository)
	mockAuthClient := new(mocks.MockAuthGrpcClient)

	var mockQuiz domain.BaseQuiz
	err := faker.FakeData(&mockQuiz)
	assert.NoError(t, err)

	var mockUser domain.User
	err = faker.FakeData(&mockUser)
	assert.NoError(t, err)

	mockQuiz.CreatorName = mockUser.Username
	t.Run("success get quiz", func(t *testing.T) {
		mockQuizRepo.On("Get", mock.Anything, mockQuiz.ID.Hex()).Return(mockQuiz, nil).Once()
		quizService := service.NewQuizService(mockQuizRepo, mockAuthClient)
		for i := 0; i < len(mockQuiz.Questions); i++ {
			mockQuiz.Questions[i].UserAnswers = []domain.UserAnswer{}
		}

		mockAuthClient.On("GetUserByID", mock.Anything, mockQuiz.CreatorID).Return(mockUser, nil).Once()

		quiz, err := quizService.Get(context.TODO(), mockQuiz.ID.Hex())

		assert.NoError(t, err)
		assert.NotEmpty(t, quiz)
		assert.Equal(t, quiz, mockQuiz)

		mockQuizRepo.AssertExpectations(t)
		mockAuthClient.AssertExpectations(t)
	})

	t.Run("not found quiz", func(t *testing.T) {
		mockQuizRepo.On("Get", mock.Anything, mockQuiz.ID.Hex()).Return(domain.BaseQuiz{}, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, fmt.Sprintf("quiz not found"))).Once()
		quizService := service.NewQuizService(mockQuizRepo, mockAuthClient)
		for i := 0; i < len(mockQuiz.Questions); i++ {
			mockQuiz.Questions[i].UserAnswers = []domain.UserAnswer{}
		}

		mockAuthClient.On("GetUserByID", mock.Anything, mockQuiz.CreatorID).Return(mockUser, nil).Once()

		quiz, err := quizService.Get(context.TODO(), mockQuiz.ID.Hex())

		assert.Error(t, err)
		assert.Equal(t, err, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, fmt.Sprintf("quiz not found")))
		assert.Empty(t, quiz)

		mockQuizRepo.AssertExpectations(t)

	})

}

func TestGetQuizByCreatorID(t *testing.T) {
	mockQuizRepo := new(mocks.MockQuizRepository)
	mockAuthClient := new(mocks.MockAuthGrpcClient)
	mockLimit := 10
	mockOffset := 0
	mockListQuiz := make([]domain.BaseQuiz, 0)

	for i := 0; i < 10; i++ {
		var newMockFakerQuiz domain.BaseQuiz
		err := faker.FakeData(&newMockFakerQuiz)
		assert.NoError(t, err)
		mockListQuiz = append(mockListQuiz, newMockFakerQuiz)
	}

	mockListUser := make([]domain.User, 0)
	for i := 0; i < 10; i++ {
		var newMockFakerUser domain.User
		err := faker.FakeData(&newMockFakerUser)
		assert.NoError(t, err)
		mockListUser = append(mockListUser, newMockFakerUser)
	}

	var mockCreatorID string = "1234567890"

	t.Run("success get quiz by creatorID", func(t *testing.T) {

		mockQuizRepo.On("GetAllQuizByCreatorID", mock.Anything, mockCreatorID, uint64(mockLimit), uint64(mockOffset)).Return(mockListQuiz, nil).Once()

		quizService := service.NewQuizService(mockQuizRepo, mockAuthClient)

		var userIDs []string
		for i := 0; i < len(mockListQuiz); i++ {

			userIDs = append(userIDs, mockListQuiz[i].CreatorID)
		}

		// get user detail
		mockAuthClient.On("GetUsersByIds", mock.Anything, userIDs).Return(mockListUser, nil).Once()

		var userIDmap map[string]string = make(map[string]string)
		for i := 0; i < len(mockListUser); i++ {
			userIDmap[mockListUser[i].ID] = mockListUser[i].Username
		}

		for i := 0; i < len(mockListQuiz); i++ {
			mockListQuiz[i].CreatorName = userIDmap[mockListQuiz[i].CreatorID]
		}

		quizs, err := quizService.GetQuizByCreatorID(context.TODO(), mockCreatorID, uint64(mockLimit), uint64(mockOffset))

		assert.NoError(t, err)
		assert.NotEmpty(t, quizs)
		assert.Equal(t, quizs, mockListQuiz)
		assert.Len(t, quizs, len(mockListQuiz))

		mockQuizRepo.AssertExpectations(t)
		mockAuthClient.AssertExpectations(t)
	})

	t.Run("not found get quiz by creatorID", func(t *testing.T) {

		mockQuizRepo.On("GetAllQuizByCreatorID", mock.Anything, mockCreatorID, uint64(mockLimit), uint64(mockOffset)).Return([]domain.BaseQuiz{}, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, fmt.Sprintf("quiz not found"))).Once()

		quizService := service.NewQuizService(mockQuizRepo, mockAuthClient)

		var userIDs []string
		for i := 0; i < len(mockListQuiz); i++ {

			userIDs = append(userIDs, mockListQuiz[i].CreatorID)
		}

		// get user detail
		mockAuthClient.On("GetUsersByIds", mock.Anything, userIDs).Return(mockListUser, nil).Once()

		var userIDmap map[string]string = make(map[string]string)
		for i := 0; i < len(mockListUser); i++ {
			userIDmap[mockListUser[i].ID] = mockListUser[i].Username
		}

		for i := 0; i < len(mockListQuiz); i++ {
			mockListQuiz[i].CreatorName = userIDmap[mockListQuiz[i].CreatorID]
		}

		quizs, err := quizService.GetQuizByCreatorID(context.TODO(), mockCreatorID, uint64(mockLimit), uint64(mockOffset))

		assert.Error(t, err)
		assert.Equal(t, err, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, fmt.Sprintf("quiz not found")))
		assert.Empty(t, quizs)

		mockQuizRepo.AssertExpectations(t)
	})
}

func TestGetQuizHistory(t *testing.T) {
	mockQuizRepo := new(mocks.MockQuizRepository)
	mockAuthClient := new(mocks.MockAuthGrpcClient)
	mockLimit := 10
	mockOffset := 0

	mockListQuiz := make([]domain.BaseQuizIsParticipant, 0)

	for i := 0; i < 10; i++ {
		var newMockFakerQuiz domain.BaseQuizIsParticipant
		err := faker.FakeData(&newMockFakerQuiz)
		assert.NoError(t, err)
		mockListQuiz = append(mockListQuiz, newMockFakerQuiz)
	}

	mockListUser := make([]domain.User, 0)
	for i := 0; i < 10; i++ {
		var newMockFakerUser domain.User
		err := faker.FakeData(&newMockFakerUser)
		assert.NoError(t, err)
		mockListUser = append(mockListUser, newMockFakerUser)
	}
	var participantID string = "1234567890"

	t.Run("success get quiz history", func(t *testing.T) {

		mockQuizRepo.On("GetQuizHistory", mock.Anything, participantID, uint64(mockLimit), uint64(mockOffset)).Return(mockListQuiz, nil).Once()

		quizService := service.NewQuizService(mockQuizRepo, mockAuthClient)

		var userIDs []string
		for i := 0; i < len(mockListQuiz); i++ {

			userIDs = append(userIDs, mockListQuiz[i].CreatorID)
		}

		// get user detail
		mockAuthClient.On("GetUsersByIds", mock.Anything, userIDs).Return(mockListUser, nil).Once()

		var userIDmap map[string]string = make(map[string]string)
		for i := 0; i < len(mockListUser); i++ {
			userIDmap[mockListUser[i].ID] = mockListUser[i].Username
		}

		for i := 0; i < len(mockListQuiz); i++ {
			mockListQuiz[i].CreatorName = userIDmap[mockListQuiz[i].CreatorID]
		}

		quizs, err := quizService.GetQuizHistory(context.TODO(), participantID, uint64(mockLimit), uint64(mockOffset))

		assert.NoError(t, err)
		assert.NotEmpty(t, quizs)
		assert.Equal(t, quizs, mockListQuiz)
		assert.Len(t, quizs, len(mockListQuiz))

		mockQuizRepo.AssertExpectations(t)
		mockAuthClient.AssertExpectations(t)
	})


	t.Run("not found get quiz history", func(t *testing.T) {
		
		mockQuizRepo.On("GetQuizHistory", mock.Anything, participantID, uint64(mockLimit), uint64(mockOffset)).Return([]domain.BaseQuizIsParticipant{}, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, fmt.Sprintf("quiz not found"))).Once()

		quizService := service.NewQuizService(mockQuizRepo, mockAuthClient)

		var userIDs []string
		for i := 0; i < len(mockListQuiz); i++ {

			userIDs = append(userIDs, mockListQuiz[i].CreatorID)
		}

		// get user detail
		mockAuthClient.On("GetUsersByIds", mock.Anything, userIDs).Return(mockListUser, nil).Once()

		var userIDmap map[string]string = make(map[string]string)
		for i := 0; i < len(mockListUser); i++ {
			userIDmap[mockListUser[i].ID] = mockListUser[i].Username
		}

		for i := 0; i < len(mockListQuiz); i++ {
			mockListQuiz[i].CreatorName = userIDmap[mockListQuiz[i].CreatorID]
		}

		quizs, err := quizService.GetQuizHistory(context.TODO(), participantID, uint64(mockLimit), uint64(mockOffset))

		assert.Error(t, err)
		assert.Equal(t, err, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, fmt.Sprintf("quiz not found")))
		assert.Empty(t, quizs)

		mockQuizRepo.AssertExpectations(t)
	})

}
