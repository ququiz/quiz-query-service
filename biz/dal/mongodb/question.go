package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"ququiz/lintang/quiz-query-service/biz/domain"
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
	// 		{"questions.user_answers", 0},
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
func (r *QuestionRepository) GetUserAnswerInAQuiz(ctx context.Context, quizID string, userID string) ([]domain.QuestionWithUserAnswerAggregate, error) {
	coll := r.db.Collection("base_quiz")
	quizIDObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return []domain.QuestionWithUserAnswerAggregate{}, err
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
			{"path", "$questions.user_answers"},
		}},
	}

	matchUser := bson.D{
		{"$match", bson.D{
			{"questions.user_answers.participant_id", userID},
		}},
	}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{matchQuizID, unwindQuestion, userAnswerFilter, matchUser})
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

func (r *QuestionRepository) GetQuestionByIDAndQuizID(ctx context.Context, quizID string, questionID string) (domain.Question, error) {
	coll := r.db.Collection("base_quiz")
	quizIDObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return domain.Question{}, err
	}

	questionObjectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		zap.L().Error("primitive.ObjectIDFromHex (quizIDObjectID) (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return domain.Question{}, err
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
		return domain.Question{}, err
	}

	var question domain.Question
	if err := cursor.Decode(&question); err != nil {
		zap.L().Error("cursor.All (GetUserAnswerInAQuiz) (QuestionRepository)", zap.Error(err))
		return domain.Question{}, err
	}

	return question, nil
}
