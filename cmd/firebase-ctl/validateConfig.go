package main

import (
	"context"
	"fmt"
	"github.com/roppenlabs/firebase-ctl/internal/firebase"
	"github.com/spf13/cobra"
	"log"
)

var inputDir string
var validateConfig = &cobra.Command{
	Use:   "remote-config",
	Short: "validate remote-config by performing a dry-run",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		clientStore, err := firebase.GetClientStore(ctx)
		if err != nil {
			log.Fatalf("Error while getting firebase app: %s", err.Error())
		}
		localConfig, err := clientStore.GetLocalConfig(inputDir)
		if err!= nil{
			log.Fatal("error reading config from local", err)
		}
		err = clientStore.ApplyConfig(*localConfig, true)
		if err != nil {
			fmt.Println(err.Error())
			return
		}


	},
}

func init() {
	validateCmd.AddCommand(validateConfig)
	validateConfig.PersistentFlags().StringVar(&inputDir, "input-dir", "", "Path to output directory")
	validateConfig.MarkPersistentFlagRequired("input-dir")
}

