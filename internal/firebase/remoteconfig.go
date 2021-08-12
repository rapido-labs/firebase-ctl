package firebase

import (
	"os"
	"path/filepath"

	"github.com/roppenlabs/firebase-ctl/internal/config"

	"github.com/roppenlabs/firebase-ctl/internal/utils"

	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
)

func (cs ClientStore) GetLatestRemoteConfig() (*remoteconfig.Response, error) {
	latestRemoteConfig, err := cs.RemoteConfigClient.GetRemoteConfig("")
	if err != nil {
		return nil, err
	}
	return latestRemoteConfig, err
}

func (cs ClientStore)BackupRemoteConfig(remoteConfig *remoteconfig.Response, outputDir string) []error {
	errs := []error{}
	conditions := remoteConfig.Conditions
	conditionsDirPath := filepath.Join(outputDir, config.ConditionsDir)
	cs.FsClient.MkdirAll(conditionsDirPath, 0744)

	for _, condition := range conditions {
		data, err := utils.JSONMarshal(condition)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		fileName := filepath.Join(conditionsDirPath, condition.Name+".json")
		file, err := cs.FsClient.OpenFile(fileName,os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644 )
		if err!= nil{
			errs = append(errs, err)
			continue
		}
		_, err = file.Write(data)
		if err != nil {
			errs = append(errs, err)
		}
	}

	parameters := remoteConfig.Parameters
	parametersDirPath := filepath.Join(outputDir, config.ParametersDir)
	cs.FsClient.MkdirAll(parametersDirPath, 0744)

	for key, parameter := range parameters {
		data, err := utils.JSONMarshal(parameter)
		if err != nil {
			errs = append(errs, err)
		}

		fileName := filepath.Join(parametersDirPath, key+".json")
		file, err := cs.FsClient.OpenFile(fileName,os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644 )
		if err!= nil{
			errs = append(errs, err)
			continue
		}
		_, err = file.Write(data)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

