package main

import (
	"github.com/spf13/cobra"
)

var validateOnly = true
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply config to remote",
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
