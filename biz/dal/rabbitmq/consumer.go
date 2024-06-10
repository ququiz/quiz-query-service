package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"ququiz/lintang/quiz-query-service/biz/service"

	"go.uber.org/zap"
)

type ScoringSvcConsumer struct {
	rmq *RabbitMQ
	rds service.CachedQsRepo
}

func NewScoringSvcConsumer(r *RabbitMQ, rds service.CachedQsRepo) *ScoringSvcConsumer {
	return &ScoringSvcConsumer{r, rds}
}

const ScoringSvcConsumerName = "quiz-query-consumer"

func (r *ScoringSvcConsumer) ListenAndServe() error {
	

	msgs, err := r.rmq.Channel.Consume(
		"delete-cache-queue",
		ScoringSvcConsumerName,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		zap.L().Fatal(fmt.Sprint("cant consume message from queue %s", "delete-cache-queue"))
	}

	go func() {
		for msg := range msgs {
			zap.L().Info("Received message: %s" + msg.RoutingKey)
			var msgBody DeleteCacheForSpecificQuizMesssage
			if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&msgBody); err != nil {
				zap.L().Error("gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&msgBody) (Scoring Consumer) (RabbitMQCOnsumer )", zap.Error(err))
				continue
			}

			var nack bool
			switch msg.RoutingKey {
			case "delete-cache":
				// TODO: implement delete cache questionns & delete cache leaderboard
				r.rds.DeleteCacheForSpecificQuiz(context.Background(), msgBody.QuizID)
			default:
				nack = true
			}

			if nack {
				zap.L().Info(fmt.Sprintf("NAcking message from queue %s", "delete-cache-queue"))

				_ = msg.Nack(false, nack)
			} else {
				zap.L().Info("Acking ")

				_ = msg.Ack(false)
			}

			zap.L().Info("No more messages to consume. Extiing.")

		}
	}()

	return nil

}

type DeleteCacheForSpecificQuizMesssage struct {
	QuizID string
}
