package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	// _ "github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	failOnError(err, "failed to connect to rabbit")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to ope channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed echange declare")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(
		ctx,
		"logs_direct",
		severityFrom(os.Args),
		false,
		false,
		amqp.Publishing{
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

func severityFrom(args []string) string {
	var s string
	if len(args) < 2 || args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}

	return s
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", msg, err)
	}
}
