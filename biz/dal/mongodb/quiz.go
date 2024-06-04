package mongodb

import (
	"context"
	"fmt"

	"ququiz/lintang/quiz-query-service/biz/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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
		zap.L().Error("coll.InsertMany (InsertQuizData) (QuizRepository)", zap.Error(err))
		return err
	}
	return nil
}

func (r *QuizRepository) GetAll(ctx context.Context) ([]domain.BaseQuiz, error) {
	coll := r.db.Collection("base_quiz")
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		zap.L().Error("coll.Find() (GetALlQuiz) (QuizRepoistory)", zap.Error(err))
		return []domain.BaseQuiz{}, err
	}

	var quizs []domain.BaseQuiz

	// bawah gakbisa
	if err := cursor.All(ctx, &quizs); err != nil {
		zap.L().Error("cursor.All()", zap.Error(err))
		return []domain.BaseQuiz{}, err
	}

	return quizs, nil
}

func (r *QuizRepository) Get(ctx context.Context, quizID string) (domain.BaseQuiz, error) {
	coll := r.db.Collection("base_quiz")
	quizIDObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return domain.BaseQuiz{}, err
	}

	filter := bson.D{{"_id", quizIDObjectID}}
	var quiz domain.BaseQuiz
	err = coll.FindOne(ctx, filter).Decode(&quiz)
	if err != nil {
		return domain.BaseQuiz{}, domain.WrapErrorf(err, domain.ErrNotFound, fmt.Sprintf(`quiz with id %s not found`, quizID))
	}
	return quiz, nil
}

func (r *QuizRepository) IsUserQuizParticipant(ctx context.Context, quizID string, userID string) ([]domain.BaseQuizIsParticipant, error) {
	quizIDObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.BaseQuizIsParticipant{}, err
	}
	filterQuiz := bson.D{
		{"$match", bson.D{
			{"_id", quizIDObjectID},
		}},
	}
	unwindParticipant := bson.D{
		{"$unwind", bson.D{
			{"path", "$participants"},
		}},
	}

	filterParticipant := bson.D{
		{"$match", bson.D{
			{"participants.user_id", userID},
		}},
	}

	coll := r.db.Collection("base_quiz")

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{filterQuiz, unwindParticipant, filterParticipant})
	if err != nil {
		zap.L().Error("coll.Aggregat (IsUserQuizParticipant) (QuizRepository)", zap.Error(err))
		return []domain.BaseQuizIsParticipant{}, err
	}

	var participant []domain.BaseQuizIsParticipant
	if err := cursor.All(ctx, &participant); err != nil {
		zap.L().Error("cursor.ALl()(IsUserQuizParticipant) (QuizRepository) ", zap.Error(err))
		return []domain.BaseQuizIsParticipant{}, err
	}

	return participant, nil
}

func (r *QuizRepository) GetAllQuizByCreatorID(ctx context.Context, creatorID string) ([]domain.BaseQuiz, error) {
	coll := r.db.Collection("base_quiz")
	filterQuiz := bson.D{

		{"creator_id", creatorID},
	}

	var quizs []domain.BaseQuiz
	cursor, err := coll.Find(ctx, filterQuiz)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return []domain.BaseQuiz{}, domain.WrapErrorf(err, domain.ErrNotFound, fmt.Sprintf(`you havent create any quiz`, creatorID))
		}

		zap.L().Error("cursor.ALl()(IsUserQuizParticipant) (QuizRepository) ", zap.Error(err))
		return []domain.BaseQuiz{}, err
	}

	if err := cursor.All(ctx, &quizs); err != nil {
		zap.L().Error("cursor.ALl()(IsUserQuizParticipant) (QuizRepository) ", zap.Error(err))
		return []domain.BaseQuiz{}, err
	}

	return quizs, nil
}

func (r *QuizRepository) GetQuizHistory(ctx context.Context, participantID string) ([]domain.BaseQuizIsParticipant, error) {
	unwindParticipant := bson.D{
		{"$unwind", bson.D{
			{"path", "$participants"},
		}},
	}
	filterParticipant := bson.D{
		{"$match", bson.D{
			{"participants.user_id", participantID},
		}},
	}

	coll := r.db.Collection("base_quiz")
	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{unwindParticipant, filterParticipant})
	if err != nil {
		zap.L().Error("coll.Aggregat (IsUserQuizParticipant) (QuizRepository)", zap.Error(err))
		return []domain.BaseQuizIsParticipant{}, err
	}

	var quizs []domain.BaseQuizIsParticipant

	if err := cursor.All(ctx, &quizs); err != nil {
		zap.L().Error("cursor.ALl()(IsUserQuizParticipant) (QuizRepository) ", zap.Error(err))
		return []domain.BaseQuizIsParticipant{}, err
	}
	return quizs, nil
}
