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

func (r *QuizRepository) InsertQuizData(ctx context.Context, quizReqs []domain.BaseQuiz) error {
	coll := r.db.Collection("base_quiz")

	var quizs []interface{}
	for _, req := range quizReqs {
		quizs = append(quizs, req)
	}
	
	_, err := coll.InsertMany(ctx, quizs)


	if err != nil {
		zap.L().Error("coll.InsertMany (InsertQuizData) (QuizRepository)", zap.Error(err ))
		return err
	}
	return nil 
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

func (r *QuizRepository) IsUserQuizParticipant(ctx context.Context, quizID string, userID string) error {
	filter := bson.D{
		{"$unwind", bson.D{
			{"path", "$participants"},
		}},
		{"$match", bson.D{
			{"$participants._id", userID},
		}},
	}

	coll := r.db.Collection("quiz")

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{filter})
	if err != nil {
		zap.L().Error("coll.Aggregat (IsUserQuizParticipant) (QuizRepository)", zap.Error(err))
		return err
	}

	var participant domain.BaseQuiz
	if err := cursor.All(ctx, &participant); err != nil {
		zap.L().Error("cursor.ALl()(IsUserQuizParticipant) (QuizRepository) ", zap.Error(err))
		return err
	}

	return nil
}
