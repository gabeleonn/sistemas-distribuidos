package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	ampq "mom/core/amqp"
	"mom/core/constants"
	"mom/core/logger"
)

type ClientConfig struct {
	ID     int
	Promos []string
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Inicia o servico de cliente",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := buildClientArgs(cmd)
		return runClient(cfg)
	},
}

func init() {
	var (
		clientID     int
		clientPromos []string
	)

	clientCmd.Flags().IntVar(&clientID, "id", 0, "Identificacao do cliente")
	clientCmd.Flags().StringSliceVar(&clientPromos, "promos", nil, "Lista de promocoes do cliente (separado por virgula)")

	if err := clientCmd.MarkFlagRequired("id"); err != nil {
		panic(fmt.Errorf("erro ao marcar flag id como obrigatoria: %w", err))
	}

	if err := clientCmd.MarkFlagRequired("promos"); err != nil {
		panic(fmt.Errorf("erro ao marcar flag promos como obrigatoria: %w", err))
	}

	rootCmd.AddCommand(clientCmd)
}

func runClient(cfg *ClientConfig) error {
	log := logger.Init(constants.ServiceClient)
	log.Printf("iniciando cliente com interesse nas promocoes: %v", cfg.Promos)

	amqpClient := ampq.New()

	defer amqpClient.Close()

	queue, err := amqpClient.DeclareQueue("", false, true, true)

	if err != nil {
		return fmt.Errorf("erro ao declarar fila do cliente: %w", err)
	}

	for _, promo := range cfg.Promos {
		if err := amqpClient.BindQueue(queue.Name, promo); err != nil {
			return fmt.Errorf("erro ao associar fila com exchange: %w", err)
		}
	}

	deliveries, err := amqpClient.Consume(
		queue.Name,
		fmt.Sprintf("%s-%d", constants.ServiceClient, cfg.ID),
	)

	if err != nil {
		return fmt.Errorf("erro ao iniciar consumo do cliente: %w", err)
	}

	log.Println("cliente consumidor ativo")

	for delivery := range deliveries {
		log.Printf("promocao recebida [%s]: %s", delivery.RoutingKey, string(delivery.Body))
	}

	return nil
}

func buildClientArgs(cmd *cobra.Command) *ClientConfig {
	id, _ := cmd.Flags().GetInt("id")
	promos, _ := cmd.Flags().GetStringSlice("promos")

	if len(promos) == 0 {
		promos = append(promos, constants.EventPromotionFeatured)
	}

	return &ClientConfig{
		ID:     id,
		Promos: promos,
	}
}
