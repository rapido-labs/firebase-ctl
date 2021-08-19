package main

import (
	"context"
	"log"

	"github.com/roppenlabs/firebase-ctl/internal/firebase"

	"github.com/spf13/cobra"
)

var outputDir string

var getRemoteConfigCmd = &cobra.Command{
	Use:   "remote-config",
	Short: "backup remote-config resources from Firebase project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		clientStore, err := firebase.GetClientStore(ctx)
		if err != nil {
			log.Fatalf("Error while getting firebase app: %s", err.Error())
		}
		latestRemoteConfig, err := clientStore.GetLatestRemoteConfig()
		if err != nil {
			log.Fatal(err)
		}
		err = clientStore.BackupRemoteConfig(latestRemoteConfig, outputDir)
		if err!= nil{
			log.Fatal(err)
		}
	},
}

func init() {
	getCmd.AddCommand(getRemoteConfigCmd)
	getRemoteConfigCmd.PersistentFlags().StringVar(&outputDir, "output-dir", "", "Path to output directory")
	getRemoteConfigCmd.MarkPersistentFlagRequired("output-dir")
}
