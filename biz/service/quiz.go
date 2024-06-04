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

type QuizService struct {
	quizRepo QuizRepository
}

func NewQuizService(qRepo QuizRepository) *QuizService {
	return &QuizService{
		qRepo,
	}
}

func (s *QuizService) GetAll(ctx context.Context) ([]domain.BaseQuiz, error) {
	quizs, err := s.quizRepo.GetAll(ctx)
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
	return quiz, nil
}

func (s *QuizService) GetQuizByCreatorID(ctx context.Context, creatorID string) ([]domain.BaseQuiz, error) {
	quizs, err := s.quizRepo.GetAllQuizByCreatorID(ctx, creatorID)
	if err != nil {
		zap.L().Error(" s.quizRepo.GetAllQuizByCreatorID (GetQuizByCreatorID ) (QuizService) ")
		return []domain.BaseQuiz{}, err
	}

	return quizs, nil
}


func (s *QuizService) GetQuizHistory(ctx context.Context, participantID string) ([]domain.BaseQuizIsParticipant, error){
	quizHistory, err := s.quizRepo.GetQuizHistory(ctx, participantID)
	if err != nil {
		zap.L().Error(" s.quizRepo.GetQuizHistory (GetQuizHistory ) (QuizService) ")
		return []domain.BaseQuizIsParticipant{}, err
	}
	return quizHistory, nil 
}