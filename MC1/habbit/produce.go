package habbit

import (
	"encoding/json"
	"fmt"

	"api-go/models"

	"github.com/rabbitmq/amqp091-go"
)

func PublishVote(vote models.Vote) error {

	if Channel == nil || Channel.IsClosed() {
		if err := reconnect(); err != nil {
			return fmt.Errorf("erro ao reconectar: %w", err)
		}
	}

	exchange := "votes_exchange"

	err := Channel.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// fila 1
	q1, err := Channel.QueueDeclare(
		"votes_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = Channel.QueueBind(
		q1.Name,
		"",
		exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// fila 2
	q2, err := Channel.QueueDeclare(
		"votes_count",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = Channel.QueueBind(
		q2.Name,
		"",
		exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(vote)
	if err != nil {
		return err
	}

	err = Channel.Publish(
		exchange,
		"",
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         body,
		},
	)

	return err
}
