package firebase

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/roppenlabs/firebase-ctl/internal/config"

	"github.com/roppenlabs/firebase-ctl/internal/utils"

	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
)

func (rc RemoteConfigClient) GetLatestRemoteConfig(ctx context.Context) (*remoteconfig.Response, error) {
	latestRemoteConfig, err := rc.Client.GetRemoteConfig("")
	if err != nil {
		return nil, err
	}
	return latestRemoteConfig, err
}

func BackupRemoteConfig(remoteConfig *remoteconfig.Response, outputDir string) []error {
	errs := []error{}

	conditions := remoteConfig.Conditions
	conditionsDirPath := getDirPath(outputDir, config.ConditionsDir)
	os.MkdirAll(conditionsDirPath, 0777)

	for _, condition := range conditions {
		data, err := utils.JSONMarshal(condition)
		if err != nil {
			errs = append(errs, err)
		}

		fileName := getFileName(conditionsDirPath, condition.Name)
		err = ioutil.WriteFile(fileName, data, 0777)
		if err != nil {
			errs = append(errs, err)
		}
	}

	parameters := remoteConfig.Parameters
	parametersDirPath := getDirPath(outputDir, config.ParametersDir)
	os.MkdirAll(parametersDirPath, 0777)

	for key, parameter := range parameters {
		data, err := utils.JSONMarshal(parameter)
		if err != nil {
			errs = append(errs, err)
		}

		fileName := getFileName(parametersDirPath, key)
		err = ioutil.WriteFile(fileName, data, 0777)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func getDirPath(parentDir, childDir string) string {
	return fmt.Sprintf("%s/%s", parentDir, childDir)
}

func getFileName(dirName, fileName string) string {
	return fmt.Sprintf("%s/%s.json", dirName, fileName)
}
