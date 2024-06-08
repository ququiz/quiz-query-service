package domain

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type QuizStatus string

// const (
// 	NOTSTARTED QuizStatus = "NOT_STARTED"
// 	INPROGRESS QuizStatus = "IN_PROGRESS"
// 	DONE       QuizStatus = "DONE"
// )

// var GetQuizStatus = map[string]QuizStatus{
// 	"NOT_STARTED": NOTSTARTED,
// 	"IN_PROGRESS": INPROGRESS,
// 	"DONE":        DONE,
// }

// type BaseQuiz struct {
// 	ID           primitive.ObjectID   `bson:"_id" json:"id"`P
// 	Name         string               `json:"name" bson:"name"`
// 	CreatorID    string               `json:"creator_id" bson:"creator_id"`
// 	Passcode     string               `json:"passcode" bson:"passcode"`
// 	StartTime    time.Time            `json:"start_time" bson:"start_time"`
// 	EndTime      time.Time            `json:"end_time" bson:"end_time"`
// 	Questions    []primitive.ObjectID `json:"questions" bson:"questions"`
// 	Participants []primitive.ObjectID `json:"participants"  bson:"participants"`
// 	Status       QuizStatus           `json:"quiz_status" bson:"quiz_status"`
// }

// type Participant struct {
// 	ID         primitive.ObjectID `bson:"_id" json:"id"`
// 	UserID     string             `json:"user_id" bson:"user_id"`
// 	FinalScore int64              `json:"final_score" bson:"final_score"`
// 	Status     QuizStatus         `json:"status" bson:"status"`
// }
// type QuestionType string

// const (
// 	MULTIPLE QuestionType = "MULTIPLE"
// 	ESSAY    QuestionType = "ESSAY"
// )

// type Choice struct {
// 	ID        primitive.ObjectID `bson:"_id" json:"id"`
// 	Text      string             `json:"text" bson:"text"`
// 	IsCorrect bool               `json:"is_correct" bson:"is_correct"`
// }

// type UserAnswer struct {
// 	ChoiceID      primitive.ObjectID `bson:"choice_id" json:"choice_id"`
// 	ParticipantID primitive.ObjectID `bson:"participant_id" json:"participant_id"`
// 	Answer        string             `json:"answer" bson:"answer"`
// }

// type Question struct {
// 	ID          primitive.ObjectID `bson:"_id" json:"id"`
// 	Question    string             `json:"question" bson:"question"`
// 	Type        QuestionType       `json:"type" bson:"type"`
// 	Choices     []Choice           `json:"choices" bson:"choices"`
// 	Weight      int32              `json:"weight" bson:"weight"`
// 	UserAnswers []UserAnswer       `json:"user_answer" bson:"user_answer"`
// }

// // -----------------------------------------------
// // yang dibawah bukan disimpen di mongodb
// // hasil aggregate get user answer

// type BaseQuizWithQuestionAggregate struct {
// 	ID           primitive.ObjectID   `bson:"_id" json:"id"`
// 	Name         string               `json:"name" bson:"name"`
// 	CreatorID    string               `json:"creator_id" bson:"creator_id"`
// 	Passcode     string               `json:"passcode" bson:"passcode"`
// 	StartTime    time.Time            `json:"start_time" bson:"start_time"`
// 	EndTime      time.Time            `json:"end_time" bson:"end_time"`
// 	Questions    []Question             `json:"questions" bson:"questions"`
// 	Participants []primitive.ObjectID `json:"participants"  bson:"participants"`
// 	Status       QuizStatus           `json:"quiz_status" bson:"quiz_status"`
// }

// type QuestionWithUserAnswerAggregate struct {
// 	ID         primitive.ObjectID `bson:"_id" json:"id"`
// 	Question   string             `json:"question" bson:"question"`
// 	Type       QuestionType       `json:"type" bson:"type"`
// 	Choices    []Choice           `json:"choices" bson:"choices"`
// 	Weight     int32              `json:"weight" bson:"weight"`
// 	UserAnswer UserAnswer         `json:"user_answer" bson:"user_answer"`
// }

// // ini custom bukan di mongodb
// type QuestionAndUserAnswer struct {
// 	Question   string             `json:"question"`
// 	Type       QuestionType       `json:"type"`
// 	Choices    []Choice           `json:"choices"`
// 	Weight     int32              `json:"weight"`
// 	UserChoice primitive.ObjectID `json:"user_choice"`
// 	UserAnswer string             `json:"user_answer"`
// }
