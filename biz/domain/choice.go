package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Choice struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Text      string             `json:"text" bson:"text"`
	IsCorrect bool               `json:"is_correct" bson:"is_correct"`
}
