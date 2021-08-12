package main

import (
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "show diff of resource changes",
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
