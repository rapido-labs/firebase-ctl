package main

import (
	"context"
	"github.com/rapido-labs/firebase-ctl/internal/utils"
	"log"

	"github.com/rapido-labs/firebase-ctl/internal/firebase"

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
			log.Fatalf("%serror while getting firebase app: %s%s", utils.Red, err.Error(), utils.Reset)
		}
		latestRemoteConfig, err := clientStore.GetLatestRemoteConfig()
		if err != nil {
			log.Fatalf("%serror getting latest remote config: %s%s", utils.Red, err.Error(), utils.Reset)
		}
		err = clientStore.BackupRemoteConfig(latestRemoteConfig, outputDir)
		if err!= nil{
			log.Fatalf("%serror backing up remote config: %s%s", utils.Red, err.Error(), utils.Reset)
		}
		log.Printf("%ssuccessfully backed up the config to %s%s", utils.Green, outputDir, utils.Reset)

	},
}

func init() {
	getCmd.AddCommand(getRemoteConfigCmd)
	getRemoteConfigCmd.PersistentFlags().StringVar(&outputDir, "output-dir", "", "Path to output directory")
	getRemoteConfigCmd.MarkPersistentFlagRequired("output-dir")
}
