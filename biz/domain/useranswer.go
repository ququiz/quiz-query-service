package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionUserAnswer struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Question      string             `json:"question" bson:"question"`
	Type          QuestionType       `json:"type" bson:"type"`
	Choices       []Choice           `json:"choices" bson:"choices"`
	Weight        int32              `json:"weight" bson:"weight"`
	CorrectAnswer string             `json:"correct_answer" bson:"correct_answer"`
	UserAnswers   UserAnswer         `json:"user_answers" bson:"user_answers"`
}

type QuizUserAnswer struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `json:"name" bson:"name"`
	CreatorID    string             `json:"creator_id" bson:"creator_id"`
	Passcode     string             `json:"passcode" bson:"passcode"`
	StartTime    time.Time          `json:"start_time" bson:"start_time"`
	EndTime      time.Time          `json:"end_time" bson:"end_time"`
	Questions    QuestionUserAnswer `json:"questions" bson:"questions"`
	Participants []Participant      `json:"participants"  bson:"participants"`
	Status       QuizStatus         `json:"quiz_status" bson:"status"`
}




