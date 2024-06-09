package service

import (
	"context"
	"fmt"
	"time"

	"ququiz/lintang/quiz-query-service/biz/domain"

	"go.uber.org/zap"
)

type QuestionRepository interface {
	GetAllByQuiz(ctx context.Context, quizID string) ([]domain.BaseQuizWithQuestionAggregate, error)
	GetUserAnswerInAQuiz(ctx context.Context, quizID string, userID string) ([]domain.QuestionUserAnswer, error)
	Get(ctx context.Context, questionID string) (domain.Question, error)
	GetQuestionByIDAndQuizID(ctx context.Context, quizID string, questionID string) (domain.BaseQuizWithOneQuestionAggregate, error)
	IsUserAnswerCorrect(ctx context.Context, quizID string, questionID string,
		userChoiceID string, userEssayAnswer string) (bool, domain.CorrectAnswer, error)
	IsUserAlreadyAnswerThisQuizID(ctx context.Context, quizID string,
		questionID string,
		userID string) (bool, error)
}

type CachedQsRepo interface {
	GetCachedQuestion(ctx context.Context, quizID string) ([]domain.Question, error)
	SetCachedQuestion(ctx context.Context, quizID string, qs []domain.Question) error
	DeleteCacheForSpecificQuiz(ctx context.Context, quizID string) error
}

type ScoringSvcProducerMQ interface {
	SendCorrectAnswer(ctx context.Context, correctAnswerMsg domain.CorrectAnswer) error
}

type QuizCommandProducerMQ interface {
	SendCorrectAnswerToQuizCommandService(ctx context.Context, userAnswerMsg domain.UserAnswerMQ) error
}

type QuestionService struct {
	questionRepo          QuestionRepository
	cachedQsRepo          CachedQsRepo
	quizRepo              QuizRepository
	scoringProducerMQ     ScoringSvcProducerMQ
	quizCommandProducerMQ QuizCommandProducerMQ
}

func NewQuestionService(questionRepo QuestionRepository, cachedQsRepo CachedQsRepo, quizRepo QuizRepository,
	sProd ScoringSvcProducerMQ, quizCommandProd QuizCommandProducerMQ) *QuestionService {
	return &QuestionService{questionRepo: questionRepo, cachedQsRepo: cachedQsRepo, quizRepo: quizRepo,
		scoringProducerMQ: sProd, quizCommandProducerMQ: quizCommandProd}
}

func (s *QuestionService) GetAllByQuiz(ctx context.Context, quizID string, userID string) ([]domain.Question, error) {
	participants, err := s.quizRepo.IsUserQuizParticipant(ctx, quizID, userID)
	if err != nil {
		return []domain.Question{}, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	if len(participants) == 0 {
		return []domain.Question{}, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("maaf anda bukan participant dari quiz ini :)"))
	}

	quiz, err := s.quizRepo.Get(ctx, quizID)

	if err != nil {
		return []domain.Question{}, domain.WrapErrorf(err, domain.ErrNotFound, fmt.Sprintf("quiz with id %s not found", quizID))
	}

	// cek apakah quiz sudah dimulai
	if time.Now().Sub(quiz.StartTime) < 0 || quiz.Status == domain.NOTSTARTED {
		return []domain.Question{}, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("quiz %s belum dimulai", quizID))
	}

	// check apakah user masih allow to liat question quiznya (time.now < quiz.endTime)
	now := time.Now()
	quizEndTime := quiz.EndTime
	if now.Unix() > quizEndTime.Unix() {
		zap.L().Debug(fmt.Sprintf("user %s not allowed to access quiz %s quiz sudah selesai", userID, quizID))
		return []domain.Question{}, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("user %s not allowed to access quiz %s quiz sudah selesai", userID, quizID))
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

func (s *QuestionService) GetUserAnswers(ctx context.Context, quizID string, userID string) ([]domain.QuestionUserAnswer, error) {
	participants, err := s.quizRepo.IsUserQuizParticipant(ctx, quizID, userID)
	if err != nil {
		return []domain.QuestionUserAnswer{}, err
	}
	if len(participants) == 0 {
		return []domain.QuestionUserAnswer{}, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("user %s not registered in quiz %s", userID, quizID))
	}

	userAnswers, err := s.questionRepo.GetUserAnswerInAQuiz(ctx, quizID, userID)
	if err != nil {
		return []domain.QuestionUserAnswer{}, err
	}

	return userAnswers, nil
}

func (s *QuestionService) UserAnswerAQuestion(ctx context.Context, quizID string, questionID string,
	userChoiceID string, userEssayAnswer string, userID string, username string) (bool, error) {

	participants, err := s.quizRepo.IsUserQuizParticipant(ctx, quizID, userID)
	if err != nil {
		return false, err
	}

	if len(participants) == 0 {
		return false, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("maaf anda bukan participant dari quiz ini"))
	}

	// cek apakah user pernah jawab pertanyaan quiz ini
	userAlreadyAnswer, err := s.questionRepo.IsUserAlreadyAnswerThisQuizID(ctx, quizID, questionID, userID)
	if userAlreadyAnswer {
		return false, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("kamu sebelumnya pernah menjawab pertanyaan ini"))
	}

	isCorrect, correctAnswer, err := s.questionRepo.IsUserAnswerCorrect(ctx, quizID, questionID, userChoiceID, userEssayAnswer)
	if err != nil {

		return false, err
	}

	correctAnswer.UserID = userID
	correctAnswer.Username = username
	if isCorrect {
		// jika jawaban user benar, maka send message to scoring service, buat dicalculate new score nya (ditambah scorenya)
		err := s.scoringProducerMQ.SendCorrectAnswer(ctx, correctAnswer)
		if err != nil {
			zap.L().Error("s.scoringProducerMQ.SendCorrectAnswer", zap.Error(err))
			return false, err
		}
	} else {
		correctAnswer.Weight = 0 // kalau jawaban salah tambah skor sebelumnya dengan angka 0 
		err := s.scoringProducerMQ.SendCorrectAnswer(ctx, correctAnswer)
		if err != nil {
			zap.L().Error("s.scoringProducerMQ.SendCorrectAnswer", zap.Error(err))
			return false, err
		}
	}
	// jika jawaban salah ya gak usah kirim ke scoring service, karan skor user akan sama (gak ditambah sama sekali)..

	// send mesage to quiz-command-service mau jawaban user benar/salah, buat insert jawaban user ke database
	err = s.quizCommandProducerMQ.SendCorrectAnswerToQuizCommandService(ctx, domain.UserAnswerMQ{
		ChoiceID:      userChoiceID,
		ParticipantID: userID,
		Answer:        userEssayAnswer,
		QuizID:        quizID,
		QuestionID:    questionID,
	})

	if err != nil {
		zap.L().Error("s.quizCommandProducerMQ.SendCorrectAnswerToQuizCommandService", zap.Error(err))
		return false, err
	}
	return isCorrect, nil
}
