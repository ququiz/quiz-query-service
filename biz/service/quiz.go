package service

import (
	"context"

	"ququiz/lintang/quiz-query-service/biz/domain"

	"go.uber.org/zap"
)

type QuizRepository interface {
	GetAll(ctx context.Context) ([]domain.BaseQuiz, error)
	IsUserQuizParticipant(ctx context.Context, quizID string, userID string) ([]domain.BaseQuizIsParticipant, error)
	GetAllQuizByCreatorID(ctx context.Context, creatorID string) ([]domain.BaseQuiz, error)
	Get(ctx context.Context, quizID string) (domain.BaseQuiz, error)
	GetQuizHistory(ctx context.Context, participantID string) ([]domain.BaseQuizIsParticipant, error)
}

type AuthGRPCClient interface {
	GetUsersByIds(ctx context.Context, userIDs []string) ([]domain.User, error)
	GetUserByID(ctx context.Context, userID string) (domain.User, error)
}

type QuizService struct {
	quizRepo   QuizRepository
	authClient AuthGRPCClient
}

func NewQuizService(qRepo QuizRepository, a AuthGRPCClient) *QuizService {
	return &QuizService{
		qRepo,
		a,
	}
}

func (s *QuizService) GetAll(ctx context.Context) ([]domain.BaseQuiz, error) {
	quizs, err := s.quizRepo.GetAll(ctx)
	var userIDs []string
	for i := 0; i < len(quizs); i++ {
		userIDs = append(userIDs, quizs[i].CreatorID)
	}
	users, err := s.authClient.GetUsersByIds(ctx, userIDs)
	if err != nil {
		return []domain.BaseQuiz{}, err
	}

	var userIDmap map[string]string = make(map[string]string)
	for i := 0; i < len(users); i++ {
		userIDmap[users[i].ID] = users[i].Username
	}

	for i := 0; i < len(quizs); i++ {
		quizs[i].CreatorName = userIDmap[quizs[i].CreatorID]
	}

	if err != nil {
		return []domain.BaseQuiz{}, err
	}

	return quizs, nil
}

func (s *QuizService) Get(ctx context.Context, quizID string) (domain.BaseQuiz, error) {
	quiz, err := s.quizRepo.Get(ctx, quizID)
	if err != nil {
		return domain.BaseQuiz{}, err
	}
	for i := 0; i < len(quiz.Questions); i++ {
		quiz.Questions[i].UserAnswers = []domain.UserAnswer{}
	}
	user, err := s.authClient.GetUserByID(ctx, quiz.CreatorID)
	if err != nil {
		zap.L().Error(" s.authClient.GetUserByID (Get) (QuizService)", zap.Error(err))
		return domain.BaseQuiz{}, err
	}
	quiz.CreatorName = user.Username
	return quiz, nil
}

func (s *QuizService) GetQuizByCreatorID(ctx context.Context, creatorID string) ([]domain.BaseQuiz, error) {
	quizs, err := s.quizRepo.GetAllQuizByCreatorID(ctx, creatorID)
	if err != nil {
		zap.L().Error(" s.quizRepo.GetAllQuizByCreatorID (GetQuizByCreatorID ) (QuizService) ")
		return []domain.BaseQuiz{}, err
	}

	var userIDs []string
	for i := 0; i < len(quizs); i++ {
		userIDs = append(userIDs, quizs[i].CreatorID)
	}
	users, err := s.authClient.GetUsersByIds(ctx, userIDs)
	if err != nil {
		return []domain.BaseQuiz{}, err
	}

	var userIDmap map[string]string = make(map[string]string)
	for i := 0; i < len(users); i++ {
		userIDmap[users[i].ID] = users[i].Username
	}

	for i := 0; i < len(quizs); i++ {
		quizs[i].CreatorName = userIDmap[quizs[i].CreatorID]
	}

	if err != nil {
		return []domain.BaseQuiz{}, err
	}

	return quizs, nil
}

func (s *QuizService) GetQuizHistory(ctx context.Context, participantID string) ([]domain.BaseQuizIsParticipant, error) {
	quizHistory, err := s.quizRepo.GetQuizHistory(ctx, participantID)
	if err != nil {
		zap.L().Error(" s.quizRepo.GetQuizHistory (GetQuizHistory ) (QuizService) ")
		return []domain.BaseQuizIsParticipant{}, err
	}

	var userIDs []string
	for i := 0; i < len(quizHistory); i++ {
		userIDs = append(userIDs, quizHistory[i].CreatorID)
	}
	users, err := s.authClient.GetUsersByIds(ctx, userIDs)
	if err != nil {
		return []domain.BaseQuizIsParticipant{}, err
	}

	var userIDmap map[string]string = make(map[string]string)
	for i := 0; i < len(users); i++ {
		userIDmap[users[i].ID] = users[i].Username
	}

	for i := 0; i < len(quizHistory); i++ {
		quizHistory[i].CreatorName = userIDmap[quizHistory[i].CreatorID]
	}

	if err != nil {
		return []domain.BaseQuizIsParticipant{}, err
	}

	return quizHistory, nil
}
