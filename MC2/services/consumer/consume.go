package consumer

import (
	"api-go/metrics"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func Consumer(jobs chan<- amqp091.Delivery) {
	conn, err := amqp091.Dial("amqps://dftnstit:fgQhLZZg5X-dDBR3egLRWdBOSdMWv4xP@woodpecker.rmq.cloudamqp.com/dftnstit")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// controle de carga
	ch.Qos(10, 0, false)

	q, err := ch.QueueDeclare(
		"votes_queue",
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

	// 🔥 só repassa os dados
	for msg := range msgs {
		jobs <- msg

		metrics.QueueSize.Set(float64(len(jobs)))
	}
}
