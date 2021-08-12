package main

import (
	"context"
	"github.com/roppenlabs/firebase-ctl/internal/firebase"
	"github.com/spf13/cobra"
	"log"
)

var configDir string

var diffRemoteConfigCmd = &cobra.Command{
	Use:   "remote-config",
	Short: "backup remote-config resources from Firebase project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		clientStore, err := firebase.GetClientStore(ctx)
		if err != nil {
			log.Fatalf("Error while getting firebase app: %s", err.Error())
		}
		clientStore.GetRemoteConfigDiff(cmd.Flag("config-dir").Value.String())

	},
}

func init() {
	diffCmd.AddCommand(diffRemoteConfigCmd)
	diffRemoteConfigCmd.PersistentFlags().StringVar(&configDir, "config-dir", "", "Path to config directory")
	diffRemoteConfigCmd.MarkPersistentFlagRequired("config-dir")
}
