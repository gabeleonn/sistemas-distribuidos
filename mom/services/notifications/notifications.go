package notifications

import (
	"mom/core/amqp"
	"mom/core/logger"
)

func Start() error {
	logger.Init("notifications")

	logger.Get().Println("servico de notificacoes iniciado")

	client := amqp.New()
	defer client.Close()

	if err := HandleConsumeCategories(client); err != nil {
		return err
	}

	if err := HandleConsumeHotDeals(client); err != nil {
		return err
	}

	return nil
}