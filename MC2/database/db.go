package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectDB() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	Pool, err = pgxpool.New(ctx, "postgres://usuario:senha@db:5432/meubanco?sslmode=disable")
	if err != nil {
		log.Fatal("Erro ao criar pool:", err)
	}

	if err = Pool.Ping(ctx); err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}

	_, err = Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS votes (
			id     SERIAL PRIMARY KEY,
			nome   TEXT NOT NULL,
			numero INT,
			email  TEXT,
			votos  TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Erro ao criar tabela:", err)
	}

	log.Println("Banco conectado!")
	return Pool
}
