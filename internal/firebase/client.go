package firebase

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	firebase "github.com/rapido-labs/firebase-admin-go/v4"
	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
	"github.com/rapido-labs/firebase-ctl/internal/config"
	"github.com/rapido-labs/firebase-ctl/internal/model"
	"github.com/rapido-labs/firebase-ctl/internal/utils"
	"github.com/spf13/afero"
	"google.golang.org/api/option"
)

type ClientStore struct {
	remoteConfigClient ConfigClient
	customFs           *customFs
}

func (cs *ClientStore) isRemoteEnabled() bool {
	return cs.remoteConfigClient != nil
}
func (cs *ClientStore) GetLatestRemoteConfig() (*remoteconfig.RemoteConfig, error) {
	if !cs.isRemoteEnabled() {
		return nil, fmt.Errorf("remote client is not configured")
	}
	latestRemoteConfigResponse, err := cs.remoteConfigClient.GetRemoteConfig("")
	if err != nil {
		return nil, err
	}

	return latestRemoteConfigResponse.RemoteConfig, err
}
func (cs *ClientStore) BackupRemoteConfig(rc *remoteconfig.RemoteConfig, outputDir string) error {
	sourceDump := model.ConvertToSourceConfig(*rc)
	for i := range sourceDump.Parameters {
		valueType := "string"
		if sourceDump.Parameters[i].DefaultValue != nil &&
			(strings.HasPrefix(sourceDump.Parameters[i].DefaultValue.ExplicitValue, "{") ||
				strings.HasPrefix(sourceDump.Parameters[i].DefaultValue.ExplicitValue, "[")) {
			valueType = "json"
		}
		parameter := sourceDump.Parameters[i]
		parameter.ValueType = valueType
		sourceDump.Parameters[i] = parameter
	}
	conditionsFilePath := filepath.Join(outputDir, config.ConditionsDir, config.ConditionsFile)
	err := cs.customFs.WriteJsonToFile(sourceDump.Conditions, conditionsFilePath)
	if err != nil {
		return fmt.Errorf("error writing to conditions file: %v", err.Error())
	}
	parameterFilePath := filepath.Join(outputDir, config.ParametersDir, config.ParametersFile)
	err = cs.customFs.WriteJsonToFile(sourceDump.Parameters, parameterFilePath)
	if err != nil {
		return fmt.Errorf("error writing to parameter file: %s", err.Error())
	}
	return nil
}
func (cs *ClientStore) GetLocalConfig(dir string) (*model.Config, error) {
	remoteConfig := &model.Config{
		Conditions:      []model.Condition{},
		Parameters:      map[string]model.Parameter{},
		ParameterGroups: nil,
	}
	conditionsFilePath := filepath.Join(dir, config.ConditionsDir, config.ConditionsFile)
	err := cs.customFs.UnmarshalFromFile(conditionsFilePath, &(remoteConfig.Conditions))
	if err != nil {
		return remoteConfig, err
	}

	parametersDirPath := filepath.Join(dir, config.ParametersDir)
	err = cs.customFs.UnMarshalFromDir(parametersDirPath, &(remoteConfig.Parameters))
	return remoteConfig, err
}
func (cs *ClientStore) pushConfigToRemote(rc remoteconfig.RemoteConfig, validateOnly bool) error {
	if !cs.isRemoteEnabled() {
		return fmt.Errorf("remote client not implemented")
	}
	template := remoteconfig.Template{
		Conditions:      rc.Conditions,
		Parameters:      rc.Parameters,
		ParameterGroups: rc.ParameterGroups,
		Version: remoteconfig.Version{
			Description:    "",
			IsLegacy:       false,
			RollbackSource: 0,
			UpdateOrigin:   "REST_API",
			UpdateTime:     time.Now(),
			UpdateType:     "FORCED_UPDATE",
			UpdateUser:     nil,
			VersionNumber:  0,
		},
	}
	_, err := cs.remoteConfigClient.PublishTemplate(context.Background(), template, validateOnly)
	if err != nil {
		return fmt.Errorf("error publishing template: %s ", err.Error())
	}
	return nil

}
func (cs *ClientStore) ValidateOnRemote(sourceConfig model.Config) error {
	rc := sourceConfig.ToRemoteConfig()
	return cs.pushConfigToRemote(*rc, true)
}
func (cs *ClientStore) ApplyConfig(sourceConfig model.Config) error {
	rc := sourceConfig.ToRemoteConfig()
	return cs.pushConfigToRemote(*rc, false)
}
func (cs *ClientStore) GetRemoteConfigDiff(inputDir string) error {
	sourceConfig, err := cs.GetLocalConfig(inputDir)
	if err != nil {
		return err
	}
	remoteConfig, err := cs.GetLatestRemoteConfig()
	if err != nil {
		return err
	}

	convertedSourceConfig := sourceConfig.ToRemoteConfig()
	utils.PrintDiff(*convertedSourceConfig, *remoteConfig)
	return nil
}

type ConfigClient interface {
	GetRemoteConfig(versionNumber string) (*remoteconfig.Response, error)
	PublishTemplate(ctx context.Context, template remoteconfig.Template, validateOnly bool) (*remoteconfig.Template, error)
}

func getFirebaseApp(ctx context.Context) (*firebase.App, error) {
	firebaseConfig, err := config.GetFirebaseConfig()
	if err != nil {
		return nil, err
	}
	opts := option.WithCredentialsFile(firebaseConfig.Service_account_json_path)

	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return nil, fmt.Errorf("error while getting firebase app: %s", err)
	}

	return app, nil
}

func GetClientStore(ctx context.Context) (*ClientStore, error) {
	firebaseApp, err := getFirebaseApp(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating firebase remote config app: %v", err.Error())
	}
	client, err := firebaseApp.RemoteConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating firebase remote config client: %v", err.Error())
	}
	return &ClientStore{remoteConfigClient: client, customFs: &customFs{afero.NewOsFs()}}, nil
}
