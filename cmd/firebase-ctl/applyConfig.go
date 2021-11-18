package main

import (
	"context"
	"github.com/rapido-labs/firebase-ctl/internal/firebase"
	"github.com/rapido-labs/firebase-ctl/internal/utils"
	"github.com/spf13/cobra"
	"log"
)

var applyConfig = &cobra.Command{
	Use:   "remote-config",
	Short: "backup remote-config resources from Firebase project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		clientStore, err := firebase.GetClientStore(ctx)
		if err != nil {
			log.Fatalf("Error while getting firebase app: %s", err.Error())
		}
		cfg, err := clientStore.GetLocalConfig(inputDir)
		if err != nil {
			log.Fatal("error getting latest config", err)
			return
		}
		err = clientStore.ApplyConfig(*cfg)
		if err != nil {
			log.Fatal("error applying latest config", err)
			return
		}
		log.Printf("%s remote config applied successfully%s", utils.Green, utils.Reset)

	},
}

func init() {
	applyCmd.AddCommand(applyConfig)
	applyConfig.PersistentFlags().StringVar(&inputDir, "input-dir", "", "Path to output directory")
	applyConfig.MarkPersistentFlagRequired("input-dir")
}
