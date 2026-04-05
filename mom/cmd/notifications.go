package cmd

import (
	"mom/services/notifications"

	"github.com/spf13/cobra"
)

var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Inicia o servico de notificacoes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notifications.Start()
	},
}

func init() {
	rootCmd.AddCommand(notificationsCmd)
}
