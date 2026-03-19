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

// 🔥 cliente HTTP otimizado
var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     30 * time.Second,
	},
	Timeout: 10 * time.Second,
}

// 🔥 limite REAL de concorrência (ESSENCIAL)
var sem = make(chan struct{}, 5) // 👈 começa baixo!

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

// 👷 worker (SEM goroutine interna)
func worker(id int, jobs <-chan Vote) {
	for v := range jobs {

		// ⏱️ comportamento humano
		time.Sleep(time.Duration(rand.Intn(1500)+300) * time.Millisecond)

		// 🔥 trava concorrência
		sem <- struct{}{}

		body, _ := json.Marshal(v)

		resp, err := client.Post(
			"https://micros-tw7a.onrender.com/vote",
			"application/json",
			bytes.NewBuffer(body),
		)

		if err != nil {
			fmt.Println("❌ Erro:", err)
			<-sem
			continue
		}

		resp.Body.Close()

		fmt.Printf("✅ [%s] %s -> %s -> %s\n",
			time.Now().Format("15:04:05"),
			v.Nome,
			v.Votos,
			resp.Status,
		)

		<-sem
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("🚀 Simulação realista iniciada...")

	jobs := make(chan Vote, 500)

	// 👥 quantidade de "usuários"
	for i := 0; i < 15; i++ {
		go worker(i, jobs)
	}

	ticker := time.NewTicker(time.Second)

	for range ticker.C {

		// 📊 tráfego variável (mais humano)
		qtd := rand.Intn(40) + 10 // 10 → 50 req/s

		// 🔥 picos ocasionais
		if rand.Intn(10) == 0 {
			qtd = rand.Intn(80) + 40
			fmt.Println("🔥 Pico de tráfego!")
		}

		fmt.Println("📊 Enviando", qtd, "requisições")

		for i := 0; i < qtd; i++ {
			jobs <- pessoa()
		}
	}
}
