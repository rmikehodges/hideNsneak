package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setup = &cobra.Command{
	Use:   "setup",
	Short: "instance parent command",
	Long:  `parent command for managing instances`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'instance --help' for usage.")
	},
}

func init() {
	rootCmd.AddCommand(setup)
}
