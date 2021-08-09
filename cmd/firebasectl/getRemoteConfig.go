package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var realm string
var configFilePath string

var getRemoteConfigCmd = &cobra.Command{
	Use:   "remoteconfig",
	Short: "get remoteconfig resources from Firebase project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("List of remote configs")
	},
}

func init() {
	getCmd.AddCommand(getRemoteConfigCmd)
}
