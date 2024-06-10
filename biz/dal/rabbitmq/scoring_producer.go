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

type ScoringServiceProducerMQ struct {
	ch *amqp.Channel
}

func NewScoringServiceProducerMQ(rmq *RabbitMQ) *ScoringServiceProducerMQ {
	return &ScoringServiceProducerMQ{
		ch: rmq.Channel,
	}
}

// jawaban user buat hitung skor user
func (s *ScoringServiceProducerMQ) SendCorrectAnswer(ctx context.Context, correctAnswerMsg domain.CorrectAnswer) error {
	return s.publish(ctx, "correct-answer", correctAnswerMsg)
}

func (s *ScoringServiceProducerMQ) publish(ctx context.Context, routingKey string, event interface{}) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		zap.L().Error("gob.NewEncoder(&b).Encode(event)", zap.Error(err))
		return err
	}

	err := s.ch.Publish(
		"scoring-quiz-query", // exchange
		routingKey,           // routing key
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
