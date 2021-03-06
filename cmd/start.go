package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	gtway "github.com/MihaiBlebea/go-checkout/gateway"
	sandbox "github.com/MihaiBlebea/go-checkout/sandbox_gateway"
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

		l := logrus.New()

		l.SetFormatter(&logrus.JSONFormatter{})
		l.SetOutput(os.Stdout)
		l.SetLevel(logrus.InfoLevel)

		gateway := gtway.New()

		sandboxGateway := sandbox.New(gateway)

		server.NewServer(sandboxGateway, l)

		return nil
	},
}
