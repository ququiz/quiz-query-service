package domain

type CorrectAnswer struct {
	Weight   uint64 `json:"weight"`
	UserID   string `json:"user_id"`
	Username string `json:"user_name"`
	QuizID   string `json:"quiz_id"`
}

type UserAnswerMQ struct {
	QuizID        string `json:"quiz_id"`
	QuestionID    string `json:"question_id"`
	ChoiceID      string `bson:"choice_id" json:"choice_id"`
	ParticipantID string `bson:"participant_id" json:"participant_id"`
	Answer        string `json:"answer" bson:"answer"`
}


