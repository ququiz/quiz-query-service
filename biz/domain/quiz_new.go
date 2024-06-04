package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizStatus string

const (
	NOTSTARTED QuizStatus = "NOT_STARTED"
	INPROGRESS QuizStatus = "IN_PROGRESS"
	DONE       QuizStatus = "DONE"
)

var GetQuizStatus = map[string]QuizStatus{
	"NOT_STARTED": NOTSTARTED,
	"IN_PROGRESS": INPROGRESS,
	"DONE":        DONE,
}

type BaseQuiz struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `json:"name" bson:"name"`
	CreatorID    string             `json:"creator_id" bson:"creator_id"`
	Passcode     string             `json:"passcode" bson:"passcode"`
	StartTime    time.Time          `json:"start_time" bson:"start_time"`
	EndTime      time.Time          `json:"end_time" bson:"end_time"`
	Questions    []Question         `json:"questions" bson:"questions"`
	Participants []Participant      `json:"participants"  bson:"participants"`
	Status       QuizStatus         `json:"quiz_status" bson:"status"`
}

type QuestionType string

const (
	MULTIPLE QuestionType = "MULTIPLE"
	ESSAY    QuestionType = "ESSAY"
)



type UserAnswer struct {
	ChoiceID      string `bson:"choice_id" json:"choice_id"`
	ParticipantID string `bson:"participant_id" json:"participant_id"`
	Answer        string             `json:"answer" bson:"answer"`
}



// -----------------------------------------------
// yang dibawah bukan disimpen di mongodb
// hasil aggregate get user answer

type BaseQuizWithQuestionAggregate struct {
	ID           primitive.ObjectID   `bson:"_id" json:"id"`
	Name         string               `json:"name" bson:"name"`
	CreatorID    string               `json:"creator_id" bson:"creator_id"`
	Passcode     string               `json:"passcode" bson:"passcode"`
	StartTime    time.Time            `json:"start_time" bson:"start_time"`
	EndTime      time.Time            `json:"end_time" bson:"end_time"`
	Questions    []Question           `json:"questions" bson:"questions"`
	Participants []Participant `json:"participants"  bson:"participants"`
	Status       QuizStatus           `json:"quiz_status" bson:"status"`
}

type QuestionWithUserAnswerAggregate struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Question   string             `json:"question" bson:"question"`
	Type       QuestionType       `json:"type" bson:"type"`
	Choices    []Choice           `json:"choices" bson:"choices"`
	Weight     int32              `json:"weight" bson:"weight"`
	UserAnswer UserAnswer         `json:"user_answers" bson:"user_answers"`
}

// ini custom bukan di mongodb
type QuestionAndUserAnswer struct {
	Question   string             `json:"question"`
	Type       QuestionType       `json:"type"`
	Choices    []Choice           `json:"choices"`
	Weight     int32              `json:"weight"`
	UserChoice string `json:"user_choice"`
	UserAnswer string             `json:"user_answer"`
}


type BaseQuizIsParticipant struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `json:"name" bson:"name"`
	CreatorID    string             `json:"creator_id" bson:"creator_id"`
	Passcode     string             `json:"passcode" bson:"passcode"`
	StartTime    time.Time          `json:"start_time" bson:"start_time"`
	EndTime      time.Time          `json:"end_time" bson:"end_time"`
	Questions    []Question         `json:"questions" bson:"questions"`
	Participants Participant      `json:"participants"  bson:"participants"`
	Status       QuizStatus         `json:"quiz_status" bson:"status"`
}
