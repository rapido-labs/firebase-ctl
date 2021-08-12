package main

import (
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate the remote config provided in input-dir",
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
