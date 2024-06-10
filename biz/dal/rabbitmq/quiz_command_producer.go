package rabbitmq

import (
	"context"
	"encoding/json"
	"ququiz/lintang/quiz-query-service/biz/domain"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type QuizCommandServiceProducerMQ struct {
	ch *amqp.Channel
}

func NewQuizCommandServiceProducerMQ(rmq *RabbitMQ) *QuizCommandServiceProducerMQ {
	return &QuizCommandServiceProducerMQ{
		ch: rmq.Channel,
	}
}

// buat simpen jaw
func (s *QuizCommandServiceProducerMQ) SendCorrectAnswerToQuizCommandService(ctx context.Context, userAnswerMsg domain.UserAnswerMQ) error {
	return s.publishToQuizCommandSvc(ctx, "user-answer", userAnswerMsg)
}

func (s *QuizCommandServiceProducerMQ) publishToQuizCommandSvc(ctx context.Context, routingKey string, event interface{}) error {

	jsonBody, err := json.Marshal(event)
	if err != nil {
		zap.L().Error("json.Marshal (publishToQuizCommandSvc) (quizCommandProducer)", zap.Error(err))
		return err
	}
	zap.L().Info("send json serialized data to quiz command service!!")
	err = s.ch.Publish(
		"quiz-command-quiz-query", // exchange
		routingKey,                // routing key
		false,
		false,
		amqp.Publishing{
			AppId:       "quiz-query-service",
			ContentType: "application/json",
			Body:        jsonBody,
			Timestamp:   time.Now(),
		})

	if err != nil {
		zap.L().Error("m.ch.Publish: ", zap.Error(err))
		return err
	}

	return nil
}
