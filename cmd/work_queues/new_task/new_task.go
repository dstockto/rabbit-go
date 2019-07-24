package main

import (
	"github.com/dstockto/rabbit-go/helper"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func main() {

	conn, err := amqp.Dial(helper.RabbitConnectionString)
	helper.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helper.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	helper.FailOnError(err, "Failed to declare a queue")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	helper.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
	log.Print("Done creating messages...")
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
