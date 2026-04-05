package cmd

import (
	"github.com/spf13/cobra"

	"mom/core/constants"
	"mom/core/logger"
	"mom/core/menu"
	"mom/core/store"
	"mom/services/gateway"
)

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Inicia o servico de gateway",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGateway()
	},
}

func init() {
	rootCmd.AddCommand(gatewayCmd)
}

func runGateway() error {
	log := logger.Init(constants.ServiceGateway)
	storage := store.NewPromotionStore()

	go func() {
		if err := gateway.PublishedPromotionsHandler(storage); err != nil {
			log.Errorf("erro ao ouvir promocoes publicadas: %v", err)
		}
	}()

	gatewayMenu := menu.NewLoop("Selecione uma opcao", []menu.Option{
		{
			Label: "Cadastrar promocao",
			Handler: func() error {
				return gateway.PublishPromotionHandler()
			},
		},
		{
			Label: "Listar promocoes publicadas",
			Handler: func() error {
				return gateway.ShowPublishedPromotionsHandler(storage)
			},
		},
		{
			Label: "Votar em promocao",
			Handler: func() error {
				return gateway.VotePromotionHandler(storage)
			},
		},
	})

	return gatewayMenu.Run()
}
