package workers

import (
	"api-go/database"
	"api-go/services/dados"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// 🔥 função exportada
func StartWorkers(jobs chan amqp091.Delivery, count int) {
	for i := 0; i < count; i++ {
		go Worker(i, jobs)
	}
}

func Worker(id int, jobs <-chan amqp091.Delivery) {
	for msg := range jobs {

		err := dados.ProcessarDados(msg.Body, database.Pool)
		if err != nil {
			log.Printf("Worker %d erro: %v\n", id, err)
			continue
		}

		msg.Ack(false)
	}
}
