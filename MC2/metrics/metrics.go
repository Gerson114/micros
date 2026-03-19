package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	JobsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_processed_total",
			Help: "Total de jobs processados",
		},
	)

	JobsErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_errors_total",
			Help: "Total de erros ao processar jobs",
		},
	)

	JobDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "job_duration_seconds",
			Help:    "Tempo de processamento dos jobs",
			Buckets: prometheus.DefBuckets,
		},
	)

	QueueSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "queue_size",
			Help: "Tamanho atual da fila interna",
		},
	)

	ActiveWorkers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_workers",
			Help: "Workers ativos processando",
		},
	)
)

func Init() {
	prometheus.MustRegister(
		JobsProcessed,
		JobsErrors,
		JobDuration,
		QueueSize,
		ActiveWorkers,
	)
}
