package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Question      string             `json:"question" bson:"question"`
	Type          QuestionType       `json:"type" bson:"type"`
	Choices       []Choice           `json:"choices" bson:"choices"`
	Weight        int32              `json:"weight" bson:"weight"`
	CorrectAnswer string             `json:"correct_answer" bson:"correct_answer"`
	UserAnswers   []UserAnswer       `json:"user_answers" bson:"user_answers"`
}






