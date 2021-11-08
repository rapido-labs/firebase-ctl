package main

import (
	"context"
	"github.com/roppenlabs/firebase-ctl/internal/firebase"
	"github.com/roppenlabs/firebase-ctl/internal/utils"
	"github.com/spf13/cobra"
	"log"
)


var diffRemoteConfigCmd = &cobra.Command{
	Use:   "remote-config",
	Short: "backup remote-config resources from Firebase project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		clientStore, err := firebase.GetClientStore(ctx)
		if err != nil {
			log.Fatalf("%serror while getting firebase app: %s%s", utils.Red, err.Error(), utils.Reset)
		}
		err = clientStore.GetRemoteConfigDiff(cmd.Flag("input-dir").Value.String())
		if err != nil {
			log.Fatalf("%serror computing diff: %s%s", utils.Red, err.Error(), utils.Reset)
		}
	},
}

func init() {
	diffCmd.AddCommand(diffRemoteConfigCmd)
	diffRemoteConfigCmd.PersistentFlags().StringVar(&inputDir, "input-dir", "", "Path to config directory")
	diffRemoteConfigCmd.MarkPersistentFlagRequired("input-dir")
}
