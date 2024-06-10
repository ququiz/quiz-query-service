package rabbitmq

import (
	"context"
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
	err = channel.Qos(
		1, 0,
		false,
	)
	if err != nil {
		zap.L().Error("err: channel.Qos" + err.Error())
	}

	// kirim jawaban user
	err = channel.ExchangeDeclare(
		"scoring-quiz-query",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("err: channel.ExchangeDeclare : " + err.Error())
	}
	// tadi sebelumnya pake direct & buat queue per replica malah mesasge yg sama diconsume 4 replica

	// buat scoring service (bikin manual di webnya aja)
	// channel.QueueDeclare(
	// 	"scoringQuizQueryQueue", // name
	// 	true,                    // durable
	// 	false,                   // delete when unused
	// 	false,                   // exclusive
	// 	false,                   // no-wait
	// 	nil,                     // arguments
	// )
	// channel.QueueBind(
	// 	"scoringQuizQueryQueue",
	// 	"correct-answer",
	// 	"scoring-quiz-query",
	// 	false,
	// 	nil,
	// )

	// kirim jawaban user
	err = channel.ExchangeDeclare(
		"quiz-command-quiz-query",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("err: channel.ExchangeDeclare : " + err.Error())
	}

	// bikin manual aja di webnya

	// // quiz command servcei
	// channel.QueueDeclare(
	// 	"userAnswerQueue", // name
	// 	true,              // durable
	// 	false,             // delete when unused
	// 	false,             // exclusive
	// 	false,             // no-wait
	// 	nil,               // arguments
	// )

	// // buat scoring ke quiz query (delete cache)
	// channel.QueueDeclare(
	// 	"delete-cache-queue",
	// 	true,  // durable
	// 	false, // delete when unused
	// 	false, // exclusive
	// 	false, // no-wait
	// 	nil,   // arguments
	// )

	// channel.QueueBind(
	// 	"delete-cache-queue",
	// 	"delete-cache",
	// 	"scoring-quiz-query",
	// 	false,
	// 	nil,
	// )

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}

}

func (r *RabbitMQ) Close(ctx context.Context) {

	r.Connection.Close()
}
