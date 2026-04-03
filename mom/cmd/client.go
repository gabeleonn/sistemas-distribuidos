package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	ampq "mom/core/amqp"
	"mom/core/constants"
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
	fmt.Printf("[Cliente %d] Iniciando cliente com interesse nas promocoes: %v\n", cfg.ID, cfg.Promos)
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
		fmt.Sprintf("%d-%d", constants.ServiceClient, cfg.ID),
	)

	if err != nil {
		return fmt.Errorf("erro ao iniciar consumo do cliente: %w", err)
	}

	fmt.Printf("[Cliente %d] Cliente consumidor ativo\n", cfg.ID)

	for delivery := range deliveries {
		fmt.Printf("[Cliente %d] Recebida promocao [%s]: %s\n", cfg.ID, delivery.RoutingKey, string(delivery.Body))
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
