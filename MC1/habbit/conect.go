package habbit

import (
	"log"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var (
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
	mu      sync.Mutex
)

func InitRabbit() error {
	var err error

	// Retry connection até 10 vezes
	for i := 0; i < 10; i++ {
		Conn, err = amqp091.Dial("amqp://admin:admin@rabbitmq:5672/")
		if err == nil {
			break
		}
		log.Printf("Tentativa %d: Aguardando RabbitMQ... %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return err
	}

	Channel, err = Conn.Channel()
	if err != nil {
		return err
	}

	// Configura QoS para evitar sobrecarga
	err = Channel.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	log.Println("✅ Conectado ao RabbitMQ com sucesso!")
	return nil
}

// Reconecta ao RabbitMQ se a conexão cair
func reconnect() error {
	mu.Lock()
	defer mu.Unlock()

	// Verifica se já está conectado
	if Channel != nil && !Channel.IsClosed() {
		return nil
	}

	log.Println("🔄 Reconectando ao RabbitMQ...")

	var err error
	for i := 0; i < 5; i++ {
		if Conn == nil || Conn.IsClosed() {
			Conn, err = amqp091.Dial("amqp://admin:admin@rabbitmq:5672/")
			if err != nil {
				log.Printf("Tentativa %d de reconexão falhou: %v", i+1, err)
				time.Sleep(time.Second)
				continue
			}
		}

		Channel, err = Conn.Channel()
		if err != nil {
			log.Printf("Erro ao criar canal: %v", err)
			time.Sleep(time.Second)
			continue
		}

		err = Channel.Qos(10, 0, false)
		if err != nil {
			return err
		}

		log.Println("✅ Reconectado com sucesso!")
		return nil
	}

	return err
}
