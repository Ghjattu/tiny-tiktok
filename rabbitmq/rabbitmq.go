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
	if RABBITMQ_USERNAME == "" {
		RABBITMQ_USERNAME = "guest"
	}

	RABBITMQ_PASSWORD := os.Getenv("RABBITMQ_PASSWORD")
	if RABBITMQ_PASSWORD == "" {
		RABBITMQ_PASSWORD = "guest"
	}

	RABBITMQ_IP := os.Getenv("RABBITMQ_IP")
	if RABBITMQ_IP == "" {
		RABBITMQ_IP = "127.0.0.1"
	}

	RABBITMQ_PORT := os.Getenv("RABBITMQ_PORT")
	if RABBITMQ_PORT == "" {
		RABBITMQ_PORT = "5672"
	}

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
