package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quill",
	Short: "Quill - a simple version control tool",
	Long:  "Quill is a lightweight version control system written in Go, inspired by Git.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand is provided, show help
		return cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
