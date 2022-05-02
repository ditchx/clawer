package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "clawer",
	}
)

// Execute runs root command
func Execute() error {
	return rootCmd.Execute()
}
