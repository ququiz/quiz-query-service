package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
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

func (s *QuizCommandServiceProducerMQ) SendCorrectAnswerToQuizCommandService(ctx context.Context, userAnswerMsg domain.UserAnswer) error {
	return s.publishToQuizCommandSvc(ctx, "user-answer", userAnswerMsg)
}

func (s *QuizCommandServiceProducerMQ) publishToQuizCommandSvc(ctx context.Context, routingKey string, event interface{}) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		zap.L().Error("gob.NewEncoder(&b).Encode(event)", zap.Error(err))
		return err
	}

	err := s.ch.Publish(
		"quiz-command-quiz-query", // exchange
		routingKey,                // routing key
		false,
		false,
		amqp.Publishing{
			AppId:       "quiz-query-service",
			ContentType: "application/x-encoding-gob",
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		})
	if err != nil {
		zap.L().Error("m.ch.Publish: ", zap.Error(err))
		return err
	}

	return nil
}
