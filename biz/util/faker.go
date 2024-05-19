package util

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"ququiz.org/lintang/quiz-query-service/biz/dal/mongodb"
	"ququiz.org/lintang/quiz-query-service/biz/domain"
	"ququiz.org/lintang/quiz-query-service/config"
)

// setiap quiz ada 20 question
// 200 quiz, 200*20 question, 200*20*4 choices
// misal 150 quiz diikuti 20 user
// dan ada 10000 users
func InsertQuizData(cfg *config.Config, mongo *mongodb.Mongodb) {
	faker := gofakeit.New(0)
	err := gofakeit.Seed(0)
	if err != nil {
		zap.L().Fatal("gofakeit.Seed()")
	}

	choicesNums := 20 * 200 * 4
	var questionChoices = make([]domain.Choice, choicesNums)
	for i := 0; i < choicesNums; i++ {
		if i%4 == 0 {

			questionChoices[i] = domain.Choice{
				ID:        primitive.NewObjectID(),
				Text:      faker.Quote(),
				IsCorrect: true,
			}
		} else {
			questionChoices[i] = domain.Choice{
				ID:        primitive.NewObjectID(),
				Text:      faker.Quote(),
				IsCorrect: false,
			}
		}
	}

	// bikin 10000 users
	var users = make([]domain.Participant, 10000)
	for i := 0; i < 10000; i++ {
		user := domain.Participant{
			ID:         primitive.NewObjectID(),
			UserID:     faker.UUID(), // ini dari postgre
			FinalScore: int64(faker.Number(0, 100)),
			Status:     domain.QuizStatus(faker.RandomString([]string{"NOT_STARTED", "IN_PROGRESS", "DONE"})),
		}
		users[i] = user
		// insert users[i] ke mongodb
		_, err := mongo.Conn.Collection("participant").InsertOne(context.Background(), user)

		if err != nil {
			zap.L().Fatal("mongo.Conn.Collection(participant).InsertOne (InsertQuizData)", zap.Error(err))
		}
	}

	var allQuizParticipants = make([]primitive.ObjectID, 20*150) // 150 quiz pertama diikuti 20 user

	var quizQuestions = make([]domain.Question, 20*200)   // 200 quiz , setiap quiz ada 20 question
	var quizParticipants = make([]domain.Participant, 20) // setiap 20 question participant harus diganti
	quizParticipants = users[0:20]

	first := 0
	participantIdx := 0
	for i := 0; i < 20*200; i++ {
		// bikin questions

		if i%20 == 0 && i < 3000 {
			randomStart := faker.Number(0, 10000-20)

			quizParticipants = users[randomStart : randomStart+20] // bikin 20 participant baru lagi, setiap quiz baru (setelah 20 question)
			// untuk 150 quiz pertama insert 20 participant setiap quiz
			loopIdx := 0
			for k := participantIdx; k < participantIdx+20; k++ {
				// fmt.Println(k)
				allQuizParticipants[k] = quizParticipants[loopIdx].ID // masukin 20 participant untuk setiap quiz
				loopIdx++
			}
			participantIdx += 20
		}

		// bikin userAnswer utk 150 quiz pertama, 150*20=3000 quizQuestions pertama
		// setiap 150 quiz pertama diikuti 20 user
		var userAnswers = make([]domain.UserAnswer, 20)
		if i < 3000 {
			choiceStartIdx := 0
			choiceEndIdx := 3

			for j := 0; j < 20; j++ {
				// jawaban user untuk question ke i
				userChoice := questionChoices[first : first+4][faker.Number(choiceStartIdx, choiceEndIdx)]
				userAnswers[j] = domain.UserAnswer{
					ChoiceID:      userChoice.ID,
					ParticipantID: quizParticipants[j].ID,
					Answer:        userChoice.Text,
				}
			}

		}

		// c[0:4], c[4:8], c[8:12], c[12:16]
		quizQuestion := domain.Question{
			ID:          primitive.NewObjectID(),
			Question:    faker.Question(),
			Type:        domain.QuestionType(faker.RandomString([]string{"MULTIPLE", "ESSAY"})),
			Choices:     questionChoices[first : first+4],
			Weight:      int32(faker.Number(1, 3)),
			UserAnswers: userAnswers,
		}
		quizQuestions[i] = quizQuestion
		first += 4

		// insert quizQuestions[i] ke mongodb
		_, err := mongo.Conn.Collection("question").InsertOne(context.Background(), quizQuestion)
		if err != nil {
			zap.L().Fatal("mongo.Conn.Collection(question).InsertOne (InsertQuizData) (util)", zap.Error(err))
		}
	}

	var quizs = make([]domain.BaseQuiz, 200)

	first = 0
	for i := 0; i < 200; i++ {
		var thisQuizParticipant []primitive.ObjectID
		var thisQuizQuestions = make([]primitive.ObjectID, 20)

		if i < 150 {
			thisQuizParticipant = allQuizParticipants[first : first+20]
		}

		thisQuestionIdx := 0
		for k := first; k < first+20; k++ {
			thisQuizQuestions[thisQuestionIdx] = quizQuestions[k].ID
			thisQuestionIdx++
		}

		quiz := domain.BaseQuiz{
			ID:           primitive.NewObjectID(),
			Name:         faker.Name(),
			CreatorID:    faker.UUID(),
			Passcode:     faker.Password(true, false, true, false, false, 5),
			StartTime:    faker.DateRange(time.Date(2023, 1, 1, 0, 0, 0, 0, &time.Location{}), time.Now()),
			EndTime:      faker.DateRange(time.Date(2023, 1, 1, 0, 0, 0, 0, &time.Location{}), time.Now()),
			Questions:    thisQuizQuestions,
			Participants: thisQuizParticipant,
			Status:       domain.QuizStatus(faker.RandomString([]string{"NOT_STARTED", "IN_PROGRESS", "DONE"})),
		}
		quizs[i] = quiz

		first += 20

		// insert quizs[i] ke mongodb
		_, err := mongo.Conn.Collection("base_quiz").InsertOne(context.Background(), quiz)
		if err != nil {
			zap.L().Error("mongo.Conn.Collection(quiz).InsertOne (InsertQuizData)", zap.Error(err))
		}
	}

}
