package main

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources from Firebase Project",
}

func init() {
	rootCmd.AddCommand(getCmd)
}
