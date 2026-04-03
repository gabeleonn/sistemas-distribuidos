package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	amqp "mom/core/amqp"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publica hello world em promocao.destaques (temporario)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amqp.New()
		defer client.Close()

		key := "promocao.destaques"
		body := []byte("hello world")

		if err := client.Publish(key, body); err != nil {
			return fmt.Errorf("erro ao publicar mensagem: %w", err)
		}

		fmt.Printf("Mensagem publicada em [%s]: %s\n", key, body)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}
