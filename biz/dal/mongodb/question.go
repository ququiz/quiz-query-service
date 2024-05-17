package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"ququiz.org/lintang/quiz-query-service/biz/domain"
)

type QuestionRepository struct {
	db *mongo.Database
}

func NewQuestionRepository(db *mongo.Database) *QuestionRepository {
	return &QuestionRepository{db}
}

func (r *QuestionRepository) GetAllByQuiz(ctx context.Context, quizID string) ([]domain.Question, error) {
	coll := r.db.Collection("quiz")
	lookup := bson.D{
		{"id", quizID},
		{"$lookup", bson.D{
			{"from", "question"},
			{"localField", "questions"},
			{"foreignField", "_id"},
			{"as", "questions"},
		}},
	}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{lookup})
	if err != nil {
		zap.L().Error("coll.Aggregrate (GetAllByQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.Question{}, err
	}

	var questions []domain.Question
	if err := cursor.All(ctx, &questions); err != nil {
		zap.L().Error("cursor.All() (GetAllByQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.Question{}, err
	}

	return questions, nil
}

// / dapetin jawaban user untuk setiap pertanyaan
func (r *QuestionRepository) GetUserAnswerInAQuiz(ctx context.Context, quizID string, userID string) ([]domain.QuestionWithUserAnswerAggregate, error) {
	coll := r.db.Collection("quiz")
	userAnswerFilter := bson.D{
		{"id", quizID},
		{"$unwind", "$user_answers"},
		{"$match", bson.D{
			{"$user_answers.participant_id", userID},
		}},
	}
	/*
		{question ... , }
	*/

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{userAnswerFilter})
	if err != nil {
		zap.L().Error("coll.Aggregate (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.QuestionWithUserAnswerAggregate{}, err
	}

	var questionsWithUserAnswer []domain.QuestionWithUserAnswerAggregate
	if err := cursor.All(ctx, &questionsWithUserAnswer); err != nil {
		zap.L().Error("cursor.All (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.QuestionWithUserAnswerAggregate{}, err
	}

	return questionsWithUserAnswer, nil 

}