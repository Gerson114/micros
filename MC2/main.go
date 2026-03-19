package main

import (
	"api-go/database"
	"api-go/metrics"
	"api-go/services/consumer"
	worker "api-go/services/workers"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabbitmq/amqp091-go"
)

func main() {

	database.ConnectDB()

	metrics.Init()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8080", nil)
	}()

	jobs := make(chan amqp091.Delivery, 10000)

	worker.StartWorkers(40, jobs)

	consumer.Consumer(jobs)
}
