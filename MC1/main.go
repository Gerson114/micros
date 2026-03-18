package main

import (
	"log"
	"net/http"
	"time"

	"api-go/habbit"
	"api-go/router"
)

func main() {

	err := habbit.InitRabbit()
	if err != nil {
		log.Fatal("Erro ao conectar RabbitMQ:", err)
	}

	router.Router()

	log.Println("Servidor rodando na porta 8080")

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
