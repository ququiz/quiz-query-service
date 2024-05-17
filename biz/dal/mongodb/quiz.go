package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"ququiz.org/lintang/quiz-query-service/biz/domain"
)

type QuizRepository struct {
	db *mongo.Database
}

func NewQuizRepository(db *mongo.Database) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) GetAll(ctx context.Context) ([]domain.BaseQuiz, error) {
	filter := bson.D{{"users_answers", 0}}
	coll := r.db.Collection("base_quiz")
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		zap.L().Error("coll.Find() (GetALlQuiz) (QuizRepoistory)", zap.Error(err))
		return []domain.BaseQuiz{}, err
	}

	var quizs []domain.BaseQuiz
	if err := cursor.All(ctx, &quizs); err != nil {
		zap.L().Error("cursor.All()", zap.Error(err))
		return []domain.BaseQuiz{}, err
	}

	return quizs, nil
}
