package main

import (
	"api-go/database"
	"api-go/services/consumer"
	"api-go/services/workers"

	"github.com/rabbitmq/amqp091-go"
)

func main() {

	database.ConnectDB()

	jobs := make(chan amqp091.Delivery, 500)

	workers.StartWorkers(jobs, 20)

	consumer.Consumer(jobs)
}
