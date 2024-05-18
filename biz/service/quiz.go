package service

import (
	"context"

	"ququiz.org/lintang/quiz-query-service/biz/domain"
)

type QuizRepository interface {
	GetAll(ctx context.Context) ([]domain.BaseQuiz, error)
	IsUserQuizParticipant(ctx context.Context, quizID string, userID string) error
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
