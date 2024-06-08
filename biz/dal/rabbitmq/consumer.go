package rabbitmq

import (
	"fmt"
	rediscache "ququiz/lintang/quiz-query-service/biz/dal/redisCache"

	"go.uber.org/zap"
)

type ScoringSvcConsumer struct {
	rmq *RabbitMQ
	rds rediscache.Redis
}

func NewScoringSvcConsumer(r *RabbitMQ, rds rediscache.Redis) *ScoringSvcConsumer {
	return &ScoringSvcConsumer{r, rds}
}

const ScoringSvcConsumerName = "quiz-query-consumer"

func (r *ScoringSvcConsumer) ListenAndServe() error {
	queue, err := r.rmq.Channel.QueueDeclare(
		"",
		false, // durable
		false, // delete when unused
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		zap.L().Fatal("cant create new queue (r.rmq.Channel.QueueDeclare) (ListenAndServe) (RMQConsumer) ", zap.Error(err))

	}
	err = r.rmq.Channel.QueueBind(
		queue.Name,
		"delete-cache",
		"scoring-quiz-query",
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("cant bind queue %s to exchange scoring-quiz-query", queue.Name))
	}
	msgs, err := r.rmq.Channel.Consume(
		queue.Name,
		ScoringSvcConsumerName,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		zap.L().Fatal(fmt.Sprint("cant consume message from queue %s", queue.Name))
	}

	go func() {
		for msg := range msgs {
			zap.L().Info("Received message: %s" + msg.RoutingKey)

			var nack bool
			switch msg.RoutingKey {
			case "delete-cache":
				// TODO: implement delete cache questionns & delete cache leaderboard

			default:
				nack = true
			}

			if nack {
				zap.L().Info(fmt.Sprintf("NAcking message from queue %s", queue.Name))

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
