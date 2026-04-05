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

	go func() {
		if err := HandleConsumeCategories(client); err != nil {
			logger.Get().Errorf("erro ao consumir categorias: %v", err)
		}
	}()

	go func() {
		if err := HandleConsumeHotDeals(client); err != nil {
			logger.Get().Errorf("erro ao consumir hot deals: %v", err)
		}
	}()

	select {}
}
