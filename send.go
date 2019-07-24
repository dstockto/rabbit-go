package main

import (
	"fmt"
	"github.com/dstockto/rabbit-go/helper"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@192.168.99.100:32769")
	helper.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helper.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	helper.FailOnError(err, "Failed to declare a queue")

	wg := sync.WaitGroup{}

	messages := 42
	wg.Add(messages)

	for i := 0; i < messages; i++ {
		go func(num int) {
			defer wg.Done()
			body := "Hello World! " + fmt.Sprint(num)
			err = ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			helper.FailOnError(err, "Failed to publish a message")
			log.Printf(" [x] Sent %s %d", body, num)
			helper.FailOnError(err, "Failed to publish a message")
		}(i)
	}

	wg.Wait()
	log.Print("Done creating messages...")
}
