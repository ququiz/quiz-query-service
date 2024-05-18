package service

import (
	"context"

	"ququiz.org/lintang/quiz-query-service/biz/domain"
)

type QuestionRepository interface {
	GetAllByQuiz(ctx context.Context, quizID string) ([]domain.Question, error)
	GetUserAnswerInAQuiz(ctx context.Context, quizID string, userID string) ([]domain.QuestionWithUserAnswerAggregate, error)
}

type CachedQsRepo interface {
	GetCachedQuestion(ctx context.Context, quizID string) ([]domain.Question, error)
	SetCachedQuestion(ctx context.Context, quizID string, qs []domain.Question) error
}

type QuestionService struct {
	questionRepo QuestionRepository
	cachedQsRepo CachedQsRepo
	QuizRepo QuizRepository
}

func NewQuestionService(questionRepo QuestionRepository, cachedQsRepo CachedQsRepo) *QuestionService {
	return &QuestionService{questionRepo: questionRepo, cachedQsRepo: cachedQsRepo}
}

func (s *QuestionService) GetAllByQuiz(ctx context.Context, quizID string, userID string) ([]domain.Question, error) {
	err := s.QuizRepo.IsUserQuizParticipant(ctx, quizID, userID)
	if err != nil {
		return []domain.Question{}, domain.WrapErrorf(err, domain.ErrUnauthorized,  "you are not authorized")
	}

	var questions []domain.Question
	questions, err = s.cachedQsRepo.GetCachedQuestion(ctx, quizID)
	if err != nil {
		// get from database
		questions, err = s.questionRepo.GetAllByQuiz(ctx, quizID)
		if err != nil {
			return []domain.Question{}, err
		}

		// set to redis
		err := s.cachedQsRepo.SetCachedQuestion(ctx, quizID, questions)
		if err != nil {
			return []domain.Question{}, err
		}
	}

	for i, _ := range questions {
		for i, _ := range questions[i].Choices {
			questions[i].Choices[i].IsCorrect = false // biar user  gak tau jawaban benernya
		}
	}

	return questions, nil
}

func (s *QuestionService) GetUserAnswers(ctx context.Context, quizID string, userID string) ([]domain.QuestionWithUserAnswerAggregate, error) {

	userAnswers, err := s.questionRepo.GetUserAnswerInAQuiz(ctx, quizID, userID)
	if err != nil {
		return []domain.QuestionWithUserAnswerAggregate{}, err 
	}

	return userAnswers, nil 
}
