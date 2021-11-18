package main

import (
	"context"
	"github.com/rapido-labs/firebase-ctl/internal/config"
	"github.com/rapido-labs/firebase-ctl/internal/firebase"
	"github.com/rapido-labs/firebase-ctl/internal/utils"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var inputDir string
var validateConfig = &cobra.Command{
	Use:   "remote-config",
	Short: "validate remote-config by performing a dry-run",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		isRemoteValidationEnabled := true
		clientStore, err := firebase.GetClientStore(ctx)
		if err != nil {
			isRemoteValidationEnabled = false
			log.Printf("%scould not find google application credentials. remote validation will not be available%s", utils.Yellow, utils.Reset)
		}
		localConfig, err := clientStore.GetLocalConfig(inputDir)
		if err != nil {
			log.Fatalf("%serror reading config from local: %s%s", utils.Red, err.Error(), utils.Reset)
		}
		errs := utils.ValidateParameters(localConfig.Parameters)
		if len(errs) != 0 {
			errStringBuilder := strings.Builder{}
			for j := range errs {
				errStringBuilder.WriteString("\n\t" + errs[j].Error())
			}
			log.Fatalf("%serror validating parameter values: %s%s", utils.Red, errStringBuilder.String(), utils.Reset)
			return
		}
		log.Printf("%sConfigValidation: Local validation successful %s", utils.Green, utils.Reset)
		if !isRemoteValidationEnabled {
			return
		}
		if _, err := config.GetFirebaseConfig(); err != nil {
			log.Printf("cannot initiate remoteConfigClient: %s", err.Error())
		}
		err = clientStore.ValidateOnRemote(*localConfig)
		if err != nil {
			log.Fatal("error validating with remote api", err)
			return
		}
		log.Printf("%sRemote validation successful %s", utils.Green, utils.Reset)

	},
}

func init() {
	validateCmd.AddCommand(validateConfig)
	validateConfig.PersistentFlags().StringVar(&inputDir, "input-dir", "", "Path to input directory")
	validateConfig.MarkPersistentFlagRequired("input-dir")
}
