package worker

import (
	"api-go/metrics"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Vote struct {
	Option string `json:"option"`
}

func StartWorkers(n int, jobs <-chan amqp091.Delivery) {
	for i := 0; i < n; i++ {
		go worker(i, jobs)
	}
}

func worker(id int, jobs <-chan amqp091.Delivery) {
	for msg := range jobs {

		start := time.Now()
		metrics.ActiveWorkers.Inc()

		var vote Vote

		err := json.Unmarshal(msg.Body, &vote)
		if err != nil {
			metrics.JobsErrors.Inc() // incrementa o contador de erros
			msg.Nack(false, false)
			continue
		}

		log.Println("✅ Worker", id, "processando:", vote)

		// 👉 aqui você salva no banco depois

		msg.Ack(false)

		//icrementar metricas
		metrics.JobsProcessed.Inc()

		//observar quanto tempo que leva para processar
		metrics.JobDuration.Observe(float64(time.Since(start).Seconds()))

		metrics.ActiveWorkers.Dec()
	}
}
