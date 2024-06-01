package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"ququiz.org/lintang/quiz-query-service/biz/domain"
)

type QuestionRepository interface {
	GetAllByQuiz(ctx context.Context, quizID string) ([]domain.BaseQuizWithQuestionAggregate, error)
	GetUserAnswerInAQuiz(ctx context.Context, quizID string, userID string) ([]domain.QuestionWithUserAnswerAggregate, error)
	Get(ctx context.Context, questionID string) (domain.Question, error)
	GetQuestionByIDAndQuizID(ctx context.Context, quizID string, questionID string) (domain.Question, error)
}

type CachedQsRepo interface {
	GetCachedQuestion(ctx context.Context, quizID string) ([]domain.Question, error)
	SetCachedQuestion(ctx context.Context, quizID string, qs []domain.Question) error
}

type QuestionService struct {
	questionRepo QuestionRepository
	cachedQsRepo CachedQsRepo
	quizRepo     QuizRepository
}

func NewQuestionService(questionRepo QuestionRepository, cachedQsRepo CachedQsRepo, quizRepo QuizRepository) *QuestionService {
	return &QuestionService{questionRepo: questionRepo, cachedQsRepo: cachedQsRepo, quizRepo: quizRepo}
}

func (s *QuestionService) GetAllByQuiz(ctx context.Context, quizID string, userID string) ([]domain.Question, error) {
	// err := s.QuizRepo.IsUserQuizParticipant(ctx, quizID, userID)
	// if err != nil {
	// 	return []domain.Question{}, domain.WrapErrorf(err, domain.ErrUnauthorized,  "you are not authorized")
	// }

	quiz, err := s.quizRepo.Get(ctx, quizID)
	// nilQuiz := domain.BaseQuiz{
	// }

	if err != nil {
		return []domain.Question{}, domain.WrapErrorf(err, domain.ErrNotFound, fmt.Sprintf("quiz with id %s not found", quizID))
	}

	// check apakah user masih allow to liat question quiznya (time.now < quiz.endTime)
	now := time.Now()
	quizEndTime := quiz.EndTime
	if now.Unix() > quizEndTime.Unix() {
		zap.L().Debug(fmt.Sprintf("user %s not allowed to access quiz %s karena waktu saat ini sudah melewati end time quiz", userID, quizID))
		return []domain.Question{}, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("user %s not allowed to access quiz %s karena waktu saat ini sudah melewati end time quiz", userID, quizID))
	}

	var questions []domain.Question
	questions, err = s.cachedQsRepo.GetCachedQuestion(ctx, quizID)
	if err != nil {
		// get from database
		quizs, err := s.questionRepo.GetAllByQuiz(ctx, quizID)
		if err != nil {
			return []domain.Question{}, err
		}

		for _, quiz := range quizs {
			questions = append(questions, quiz.Questions...)
		}
		// set to redis
		err = s.cachedQsRepo.SetCachedQuestion(ctx, quizID, questions)
		if err != nil {
			return []domain.Question{}, err
		}
	} else {
		zap.L().Debug(fmt.Sprintf("questions utk quiz %s ada di cache", quizID))
	}

	for i, _ := range questions {
		if len(questions[i].Choices) > 0 {
			for j, _ := range questions[i].Choices {
				questions[i].Choices[j].IsCorrect = false // biar user  gak tau jawaban benernya
			}
		}

	}

	return questions, nil
}

func (s *QuestionService) GetAllByQuizNotCached(ctx context.Context, quizID string, userID string) ([]domain.Question, error) {
	// err := s.QuizRepo.IsUserQuizParticipant(ctx, quizID, userID)
	// if err != nil {
	// 	return []domain.Question{}, domain.WrapErrorf(err, domain.ErrUnauthorized,  "you are not authorized")
	// }

	var questions []domain.Question

	// get from database
	quizs, err := s.questionRepo.GetAllByQuiz(ctx, quizID)
	if err != nil {
		return []domain.Question{}, err
	}

	for _, quiz := range quizs {
		questions = append(questions, quiz.Questions...)
	}

	for i, _ := range questions {
		for i, _ := range questions[i].Choices {
			questions[i].Choices[i].IsCorrect = false // biar user  gak tau jawaban benernya
		}
	}

	return questions, nil
}

func (s *QuestionService) GetUserAnswers(ctx context.Context, quizID string, userID string) ([]domain.QuestionWithUserAnswerAggregate, error) {
	participants, err := s.quizRepo.IsUserQuizParticipant(ctx, quizID, userID)
	if err != nil {
		return []domain.QuestionWithUserAnswerAggregate{}, err
	}
	if len(participants) == 0 {
		return []domain.QuestionWithUserAnswerAggregate{}, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("user %s not registered in quiz %s", userID, quizID))
	}

	userAnswers, err := s.questionRepo.GetUserAnswerInAQuiz(ctx, quizID, userID)
	if err != nil {
		return []domain.QuestionWithUserAnswerAggregate{}, err
	}

	return userAnswers, nil
}
