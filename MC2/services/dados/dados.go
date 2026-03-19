package dados

import (
	"api-go/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProcessarDados(body []byte, pool *pgxpool.Pool) error {
	var vote models.Vote

	if err := json.Unmarshal(body, &vote); err != nil {
		return err
	}

	if vote.Votos == "" || vote.Nome == "" {
		return fmt.Errorf("nome e votos obrigatorios")
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("erro ao adquirir conexao: %w", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		"INSERT INTO votes (nome, numero, email, votos) VALUES ($1, $2, $3, $4)",
		pgx.QueryExecModeSimpleProtocol,
		vote.Nome, vote.Numero, vote.Email, vote.Votos,
	)
	if err != nil {
		return fmt.Errorf("erro ao salvar voto: %w", err)
	}

	fmt.Println("Voto salvo:", vote)
	return nil
}
