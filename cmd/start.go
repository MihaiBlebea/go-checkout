package cmd

import (
	"github.com/spf13/cobra"

	gtway "github.com/MihaiBlebea/go-checkout/gateway"
	"github.com/MihaiBlebea/go-checkout/server"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application server.",
	Long:  "Start the application server.",
	RunE: func(cmd *cobra.Command, args []string) error {

		gateway := gtway.New()

		server.NewServer(gateway)

		return nil
	},
}
