package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database.",
	Long:  "Migrate the database.",
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
