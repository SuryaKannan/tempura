package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tempura",
	Short: "lightweight durable workflow engine",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nUse tempura --help to see what commands are available!")
	},
	Version: "0.1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
