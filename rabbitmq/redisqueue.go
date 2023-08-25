package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RedisMsgQueue struct {
	*RabbitMQ
	Queue     amqp.Queue
	QueueName string
}

var RedisMQ = &RedisMsgQueue{}

func InitRedisMQ(queueName string) {
	RedisMQ.RabbitMQ = MyRabbitMQ
	RedisMQ.QueueName = queueName

	var err error

	RedisMQ.Queue, err = RedisMQ.Channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	RedisMQ.failOnError("failed to declare a queue", err)

	go RedisMQ.Consumer()
}

func (rmq *RedisMsgQueue) Producer(message *Message) {
	bytes, err := json.Marshal(message)
	rmq.failOnError("failed to marshal message", err)

	err = rmq.Channel.PublishWithContext(context.Background(),
		"",
		rmq.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		},
	)
	rmq.failOnError("failed to publish a message", err)
}

func (rmq *RedisMsgQueue) Consumer() {
	msgs, err := rmq.Channel.Consume(
		rmq.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	rmq.failOnError("failed to register a consumer", err)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			message := &Message{}
			err := json.Unmarshal(d.Body, message)
			rmq.failOnError("failed to unmarshal message", err)

			ConsumeMessage(message)
		}
	}()

	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}
