package rabbitmq

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

var MyRabbitMQ = &RabbitMQ{}
var err error

func init() {
	godotenv.Load("../.env")

	RABBITMQ_USERNAME := os.Getenv("RABBITMQ_USERNAME")
	RABBITMQ_PASSWORD := os.Getenv("RABBITMQ_PASSWORD")
	RABBITMQ_IP := os.Getenv("RABBITMQ_IP")
	RABBITMQ_PORT := os.Getenv("RABBITMQ_PORT")

	url := "amqp://" + RABBITMQ_USERNAME + ":" + RABBITMQ_PASSWORD + "@" +
		RABBITMQ_IP + ":" + RABBITMQ_PORT + "/"

	var err error
	MyRabbitMQ.Conn, err = amqp.Dial(url)
	MyRabbitMQ.failOnError("failed to connect to RabbitMQ", err)

	CreateChannel()
}

func CreateChannel() {
	MyRabbitMQ.Channel, err = MyRabbitMQ.Conn.Channel()
	MyRabbitMQ.failOnError("failed to open a channel", err)

	InitRedisMQ("redis_msg_queue")
}

func (r *RabbitMQ) failOnError(msg string, err error) {
	if err != nil {
		defer r.Destroy()
		log.Panicf("%s: %s", msg, err.Error())
	}
}

func (r *RabbitMQ) Destroy() {
	r.Channel.Close()
	r.Conn.Close()
}
