package main

import (
	"encoding/json"
	"fmt"
	"log"

	"api-go/models"

	"github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	queueName := "votes_queue"

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	var vote models.Vote

	for msg := range msgs {

		err := json.Unmarshal(msg.Body, &vote)
		if err != nil {
			fmt.Println("erro ao converter:", err)
			continue
		}

		fmt.Println("Voto recebido:", vote)

		// aqui você salvaria no banco

		msg.Ack(false)
	}
}
