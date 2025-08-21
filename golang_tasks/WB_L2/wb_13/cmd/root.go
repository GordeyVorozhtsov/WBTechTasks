package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cut",
	Short: "Cut selected current column",
	Long:  `Cut is utility that can extract current columns with delimiter`,
}

func Execute() error {
	return rootCmd.Execute()
}
