package pipeline

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func InitAmqp() *amqp.Connection {
	conn, err := amqp.Dial("amqp://shareloadmin:shareloadmin@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func DeclareSendCh(conn *amqp.Connection) (amqp.Queue, *amqp.Channel) {
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	q, err := ch.QueueDeclare(
		"transcoded_queue", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	FailOnError(err, "Failed to declare a queue")
	return q, ch
}

func PublishAfterTranscoded(body []byte) {
	conn := InitAmqp()
	q, ch := DeclareSendCh(conn)
	defer ch.Close()
	defer conn.Close()
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	FailOnError(err, "Failed to publish a message")
}
