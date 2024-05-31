package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Participant struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	UserID     string             `json:"user_id" bson:"user_id"`
	FinalScore int64              `json:"final_score" bson:"final_score"`
	Status     QuizStatus         `json:"status" bson:"status"`
}
