package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		MaxConnsPerHost:     1000,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	},
}

type Voto struct {
	ID     int    `json:"id"`
	Nome   string `json:"nome"`
	Numero int    `json:"numero"`
	Email  string `json:"email"`
	Votos  string `json:"votos"`
}

// envia um voto para a API
func enviarVoto(id int) {
	voto := Voto{
		ID:     id,
		Nome:   fmt.Sprintf("Pessoa_%d", id),
		Numero: id % 100,
		Email:  fmt.Sprintf("pessoa%d@email.com", id),
		Votos:  "1",
	}

	jsonData, _ := json.Marshal(voto)

	for tentativa := 0; tentativa < 3; tentativa++ {
		resp, err := httpClient.Post(
			"http://localhost:8888/vote",
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
				return
			}
		}

		if tentativa < 2 {
			time.Sleep(time.Duration(500+rand.Intn(500)) * time.Millisecond)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	totalUsuarios := 10000
	duracaoSimulacao := 5 * time.Minute
	concorrenciaMaxima := 200

	semaphore := make(chan struct{}, concorrenciaMaxima)
	intervaloChegada := duracaoSimulacao / time.Duration(totalUsuarios)

	fmt.Printf("🚀 Simulação realista: %d usuários em %v\n", totalUsuarios, duracaoSimulacao)
	fmt.Printf("📊 Taxa: ~%.0f usuários/segundo\n", float64(totalUsuarios)/duracaoSimulacao.Seconds())
	fmt.Printf("⚡ Concorrência máxima: %d\n\n", concorrenciaMaxima)

	inicio := time.Now()
	contador := 0

	for i := 1; i <= totalUsuarios; i++ {
		semaphore <- struct{}{}

		go func(userID int) {
			defer func() { <-semaphore }()

			time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
			enviarVoto(userID)

			contador++
			if contador%500 == 0 {
				elapsed := time.Since(inicio)
				fmt.Printf("✅ %d votos enviados em %v (%.1f/s)\n",
					contador, elapsed.Round(time.Second), float64(contador)/elapsed.Seconds())
			}
		}(i)

		time.Sleep(intervaloChegada)
	}

	for i := 0; i < concorrenciaMaxima; i++ {
		semaphore <- struct{}{}
	}

	fmt.Printf("\n🎉 Simulação concluída em %v\n", time.Since(inicio).Round(time.Second))
}
