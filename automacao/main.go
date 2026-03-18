package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Vote struct {
	Nome   string `json:"nome"`
	Numero int    `json:"numero"`
	Email  string `json:"email"`
	Votos  string `json:"votos"`
}

var nomes = []string{
	"Ana Silva", "Carlos Souza", "Mariana Lima", "Pedro Oliveira",
	"Juliana Costa", "Rafael Santos", "Fernanda Rocha", "Lucas Pereira",
	"Beatriz Alves", "Thiago Martins", "Camila Ferreira", "Bruno Carvalho",
	"Larissa Gomes", "Diego Ribeiro", "Patrícia Nunes", "Rodrigo Mendes",
}

var candidatos = []string{"candidato1", "candidato2", "candidato3"}

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     30 * time.Second,
	},
	Timeout: 5 * time.Second,
}

func pessoa() Vote {
	nome := nomes[rand.Intn(len(nomes))]
	numero := rand.Intn(900) + 100
	return Vote{
		Nome:   nome,
		Numero: numero,
		Email:  fmt.Sprintf("%s%d@email.com", nome[:3], numero),
		Votos:  candidatos[rand.Intn(len(candidatos))],
	}
}

func worker(jobs <-chan Vote) {
	for v := range jobs {
		body, _ := json.Marshal(v)
		resp, err := client.Post("http://localhost:80/vote", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Erro:", err)
			continue
		}
		resp.Body.Close()
		fmt.Printf("[%s] %s -> %s -> %s\n", time.Now().Format("15:04:05"), v.Nome, v.Votos, resp.Status)
	}
}

func main() {
	fmt.Println("Iniciando simulação de votos...")

	jobs := make(chan Vote, 300)

	for i := 0; i < 50; i++ {
		go worker(jobs)
	}

	for {
		for i := 0; i < 300; i++ {
			jobs <- pessoa()
		}
		time.Sleep(1 * time.Second)
	}
}
