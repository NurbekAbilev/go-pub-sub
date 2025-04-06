package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"
	_ "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://user:password@localhost:5672/")
	failOnError(err, "failed to connect to rabbit")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to ope channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", false, false, false, false, nil,
	)
	failOnError(err, "failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	failOnError(err, "failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)

}

func bodyFrom(args []string) string {
	var s string
	if len(args) < 2 || os.Args[0] == "" {
		s = "Hello"
	} else {
		s = strings.Join(args[1:], " ")
	}

	return s
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", msg, err)
	}
}
