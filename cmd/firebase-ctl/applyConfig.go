package main

import (
	"context"
	"fmt"
	"github.com/roppenlabs/firebase-ctl/internal/firebase"
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
		err = clientStore.ApplyConfig(cmd.Flag("input-dir").Value.String(), false)
		if err != nil {
			fmt.Println(err.Error())
			return
		}


	},
}

func init() {
	applyCmd.AddCommand(applyConfig)
	applyConfig.PersistentFlags().StringVar(&outputDir, "input-dir", "", "Path to output directory")
	applyConfig.MarkPersistentFlagRequired("input-dir")
}

