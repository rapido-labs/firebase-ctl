package main

import (
	"context"
	"fmt"
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

		remoteConfigClient, err := firebase.GetRemoteConfigClient(ctx)
		if err != nil {
			log.Fatalf("Error while getting firebase app: %s", err.Error())
		}
		latestRemoteConfig, err := remoteConfigClient.GetLatestRemoteConfig(ctx)
		if err != nil {
			log.Fatal(err)
		}
		errs := firebase.BackupRemoteConfig(latestRemoteConfig, outputDir)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Println(err)
			}
			log.Fatal()
		}
	},
}

func init() {
	getCmd.AddCommand(getRemoteConfigCmd)
	getRemoteConfigCmd.PersistentFlags().StringVar(&outputDir, "output-dir", "", "Path to output directory")
	getRemoteConfigCmd.MarkPersistentFlagRequired("output-dir")
}
