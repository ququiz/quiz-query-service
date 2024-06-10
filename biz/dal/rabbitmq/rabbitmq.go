package rabbitmq

import (
	"ququiz/lintang/quiz-query-service/config"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

const quizQueryConsumerName = "quizQueryConsumer"

func NewRabbitMQ(cfg *config.Config) *RabbitMQ {
	zap.L().Info("rmq address: " + cfg.RabbitMQ.RMQAddress)

	conn, err := amqp.Dial(cfg.RabbitMQ.RMQAddress)

	if err != nil {
		zap.L().Fatal("error: cannot connect to rabbitmq: " + err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("error can't get rabbitmq cahnnel: " + err.Error())
	}

	err = channel.ExchangeDeclare(
		"scoring-quiz-query",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("err: channel.ExchangeDeclare : " + err.Error())
	}

	channel.QueueDeclare(
		"scoringQuizQueryQueue", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		zap.L().Info("err: channel.QuueeDeclare(scoringQuizQueryQueue) : " + err.Error())

	}

	channel.QueueDeclare(
		"delete-cache-queue",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	channel.QueueBind(
		"delete-cache-queue",
		"delete-cache",
		"scoring-quiz-query",
		false,
		nil,
	)

	err = channel.ExchangeDeclare(
		"quiz-command-quiz-query",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("err: channel.ExchangeDeclare : " + err.Error())
	}

	channel.QueueDeclare(
		"userAnswerQueue", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		zap.L().Info("err: channel.QuueeDeclare(userAnswerQueue) : " + err.Error())
	}

	err = channel.Qos(
		1, 0,
		false,
	)
	if err != nil {
		zap.L().Fatal("err: channel.Qos" + err.Error())
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}

}

func (r *RabbitMQ) Close() error {

	return r.Connection.Close()
}
