package mongodb

import (
	"context"
	"fmt"
	"time"

	"ququiz/lintang/quiz-query-service/biz/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type QuestionRepository struct {
	db *mongo.Database
}

func NewQuestionRepository(db *mongo.Database) *QuestionRepository {
	return &QuestionRepository{db}
}

func (r *QuestionRepository) Get(ctx context.Context, questionID string) (domain.Question, error) {
	coll := r.db.Collection("base_quiz")
	questionObjectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return domain.Question{}, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)

	}

	filterByID := bson.D{{"_id", questionObjectID}}
	cursor := coll.FindOne(ctx, filterByID)

	var question domain.Question

	if err := cursor.Decode(&question); err != nil {
		zap.L().Error("cursor.Decode (Get) (ContainerRepository), ", zap.Error(err))
		return domain.Question{}, err
	}
	return question, nil

}

func (r *QuestionRepository) GetAllByQuiz(ctx context.Context, quizID string) ([]domain.BaseQuizWithQuestionAggregate, error) {
	coll := r.db.Collection("base_quiz")
	quizIDObjectID, err := primitive.ObjectIDFromHex(quizID)

	match := bson.D{
		{"$match", bson.D{
			{"_id", quizIDObjectID},
		}},
	}
	// lookup := bson.D{

	// 	{"$lookup", bson.D{
	// 		{"from", "question"},
	// 		{"localField", "questions"},
	// 		{"foreignField", "_id"},
	// 		{"as", "questions"},
	// 	}},
	// }

	// project := bson.D{
	// 	{"$project", bson.D{
	// 		{"questions.user_answer", 0},
	// 	}},
	// }

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{match})
	if err != nil {
		zap.L().Error("coll.Aggregrate (GetAllByQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.BaseQuizWithQuestionAggregate{}, err
	}

	var questions []domain.BaseQuizWithQuestionAggregate
	if err := cursor.All(ctx, &questions); err != nil {
		zap.L().Error("cursor.All() (GetAllByQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.BaseQuizWithQuestionAggregate{}, err
	}

	return questions, nil
}

// / dapetin jawaban user untuk setiap pertanyaan
func (r *QuestionRepository) GetUserAnswerInAQuiz(ctx context.Context, quizID string, userID string) ([]domain.QuestionUserAnswer, error) {
	coll := r.db.Collection("base_quiz")
	quizIDObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.QuestionUserAnswer{}, err
	}

	matchQuizID := bson.D{
		{"$match", bson.D{
			{"_id", quizIDObjectID},
		}},
	}

	unwindQuestion := bson.D{

		{"$unwind", bson.D{
			{"path", "$questions"},
		}},
	}

	userAnswerFilter := bson.D{
		{"$unwind", bson.D{
			{"path", "$questions.user_answer"},
		}},
	}

	matchUser := bson.D{
		{"$match", bson.D{
			{"questions.user_answer.participant_id", userID},
		}},
	}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{matchQuizID, unwindQuestion, userAnswerFilter, matchUser})
	if err != nil {
		zap.L().Error("coll.Aggregate (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.QuestionUserAnswer{}, err
	}

	var questionsWithUserAnswer []domain.QuizUserAnswer
	if err := cursor.All(ctx, &questionsWithUserAnswer); err != nil {
		zap.L().Error("cursor.All (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.QuestionUserAnswer{}, err
	}

	if len(questionsWithUserAnswer) > 0 && time.Now().Sub(questionsWithUserAnswer[0].StartTime) < 0 {
		return []domain.QuestionUserAnswer{}, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("quiz %s belum dimulai :)", quizID))
	}

	var userAnswer []domain.QuestionUserAnswer
	for i := 0; i < len(questionsWithUserAnswer); i++ {
		userAnswer = append(userAnswer, domain.QuestionUserAnswer{
			ID:            questionsWithUserAnswer[i].Questions.ID,
			Question:      questionsWithUserAnswer[i].Questions.Question,
			Type:          questionsWithUserAnswer[i].Questions.Type,
			Choices:       questionsWithUserAnswer[i].Questions.Choices,
			Weight:        questionsWithUserAnswer[i].Questions.Weight,
			CorrectAnswer: questionsWithUserAnswer[i].Questions.CorrectAnswer,
			UserAnswers:   questionsWithUserAnswer[i].Questions.UserAnswers,
		})
	}

	return userAnswer, nil
}

func (r *QuestionRepository) IsUserAlreadyAnswerThisQuizID(ctx context.Context, quizID string,
	questionID string,
	userID string) (bool, error) {
	coll := r.db.Collection("base_quiz")
	quizIDObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return false, err
	}

	matchQuizID := bson.D{
		{"$match", bson.D{
			{"_id", quizIDObjectID},
		}},
	}

	unwindQuestion := bson.D{

		{"$unwind", bson.D{
			{"path", "$questions"},
		}},
	}

	questionObjectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return false, err
	}

	matchQuestion := bson.D{
		{"$match", bson.D{
			{"questions._id", questionObjectID},
		}},
	}

	userAnswerUnwind := bson.D{
		{"$unwind", bson.D{
			{"path", "$questions.user_answer"},
		}},
	}

	matchUser := bson.D{
		{"$match", bson.D{
			{"questions.user_answer.participant_id", userID},
		}},
	}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{matchQuizID, unwindQuestion, matchQuestion, userAnswerUnwind, matchUser})
	if err != nil {
		zap.L().Error("coll.Aggregate (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return false, err
	}

	var questionsWithUserAnswer []domain.QuizUserAnswer
	if err := cursor.All(ctx, &questionsWithUserAnswer); err != nil {
		zap.L().Error("cursor.All (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return false, err
	}

	if len(questionsWithUserAnswer) > 0 {
		// user pernah jawab pertanyaan ini
		return true, nil
	}
	return false, nil
}

func (r *QuestionRepository) IsUserAnswerCorrect(ctx context.Context, quizID string, questionID string,
	userChoiceID string, userEssayAnswer string) (bool, domain.CorrectAnswer, error) {
	coll := r.db.Collection("base_quiz")
	quizIDObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (IsUserAnswerCorrect) (QuestionRepository)", zap.Error(err))
		return false, domain.CorrectAnswer{}, err
	}

	questionObjectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (IsUserAnswerCorrect) (QuestionRepository)", zap.Error(err))
		return false, domain.CorrectAnswer{}, err
	}

	matchQuizID := bson.D{
		{"$match", bson.D{
			{"_id", quizIDObjectID},
		}},
	}

	unwindQuestion := bson.D{

		{"$unwind", bson.D{
			{"path", "$questions"},
		}},
	}

	matchQuestion := bson.D{
		{"$match", bson.D{
			{"questions._id", questionObjectID},
		}},
	}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{matchQuizID, unwindQuestion, matchQuestion})
	if err != nil {
		zap.L().Error("coll.Aggregate (IsUserAnswerCorrect) (QuestionRepository)", zap.Error(err))
		return false, domain.CorrectAnswer{}, err
	}

	var quiz []domain.BaseQuizWithOneQuestionAggregate

	if err := cursor.All(ctx, &quiz); err != nil {
		zap.L().Error("cursor.All (IsUserAnswerCorrect) (QuestionRepository)", zap.Error(err))
		return false, domain.CorrectAnswer{}, err
	}

	if len(quiz) == 0 {
		return false, domain.CorrectAnswer{}, domain.WrapErrorf(err, domain.ErrNotFound, fmt.Sprintf("question %s not found", questionID))
	}
	correctAnswer := domain.CorrectAnswer{
		Weight: uint64(quiz[0].Questions.Weight),
		QuizID: quizID,
	}

	// cek apakah quiz masih berjalan...
	if time.Now().Sub(quiz[0].EndTime) > 0 {
		return false, domain.CorrectAnswer{}, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("quiz %s sudah selesai", quiz[0].Name))
	}

	// cek apakah quiz sudah dimulai
	if time.Now().Sub(quiz[0].StartTime) < 0 || quiz[0].Status == domain.NOTSTARTED {
		return false, domain.CorrectAnswer{}, domain.WrapErrorf(err, domain.ErrBadParamInput, fmt.Sprintf("quiz %s belum dimulai", quiz[0].Name))

	}

	if quiz[0].Questions.Type == domain.ESSAY {
		return quiz[0].Questions.CorrectAnswer == userEssayAnswer, correctAnswer, nil
	} else {
		var correctChoiceID string = ""
		for i := 0; i < len(quiz[0].Questions.Choices); i++ {
			if quiz[0].Questions.Choices[i].IsCorrect {
				correctChoiceID = quiz[0].Questions.Choices[i].ID.Hex()
			}
		}
		return correctChoiceID == userChoiceID, correctAnswer, nil
	}

}

func (r *QuestionRepository) GetQuestionByIDAndQuizID(ctx context.Context, quizID string, questionID string) (domain.BaseQuizWithOneQuestionAggregate, error) {
	coll := r.db.Collection("base_quiz")
	quizIDObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return domain.BaseQuizWithOneQuestionAggregate{}, err
	}

	questionObjectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return domain.BaseQuizWithOneQuestionAggregate{}, err
	}

	matchQuizID := bson.D{
		{"$match", bson.D{
			{"_id", quizIDObjectID},
		}},
	}

	unwindQuestion := bson.D{

		{"$unwind", bson.D{
			{"path", "$questions"},
		}},
	}

	matchQuestion := bson.D{
		{"$match", bson.D{
			{"questions._id", questionObjectID},
		}},
	}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{matchQuizID, unwindQuestion, matchQuestion})
	if err != nil {
		zap.L().Error("coll.Aggregate (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return domain.BaseQuizWithOneQuestionAggregate{}, err
	}

	var question domain.BaseQuizWithOneQuestionAggregate
	if err := cursor.All(ctx, &question); err != nil {
		zap.L().Error("cursor.All (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return domain.BaseQuizWithOneQuestionAggregate{}, err
	}

	return question, nil
}
