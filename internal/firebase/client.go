package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	firebase "github.com/rapido-labs/firebase-admin-go/v4"
	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
	"github.com/roppenlabs/firebase-ctl/internal/config"
	"github.com/roppenlabs/firebase-ctl/internal/utils"
	"github.com/spf13/afero"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ClientStore struct {
	remoteConfigClient ConfigClient
	fsClient           afero.Fs
}

func (cs *ClientStore) GetLatestRemoteConfig() (*remoteconfig.RemoteConfig, error) {
	latestRemoteConfigResponse, err := cs.remoteConfigClient.GetRemoteConfig("")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	rc := latestRemoteConfigResponse.RemoteConfig
	for i := range rc.Parameters {
		valueType := "string"
		if strings.HasPrefix(rc.Parameters[i].DefaultValue.ExplicitValue, "{") || strings.HasPrefix(rc.Parameters[i].DefaultValue.ExplicitValue, "[") {
			valueType = "json"
		}
		parameter := rc.Parameters[i]
		parameter.ValueType = valueType
		rc.Parameters[i] = parameter
	}

	return rc, err
}
func (cs *ClientStore) BackupRemoteConfig(remoteConfig *remoteconfig.RemoteConfig, outputDir string) error {
	conditions := remoteConfig.Conditions
	conditionsDir := filepath.Join(outputDir, config.ConditionsDir)
	conditionsFilePath := filepath.Join(conditionsDir, config.ConditionsFile)
	cs.fsClient.MkdirAll(conditionsDir, 0744)
	conditionsData, err := utils.JSONMarshal(conditions)
	if err != nil {
		return fmt.Errorf("error marshalling conditions to json: %s", err.Error())
	}
	conditionsFile, err := cs.fsClient.OpenFile(conditionsFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening conditions file for write: %v", err.Error())
	}
	_, err = conditionsFile.Write(conditionsData)
	if err != nil {
		return fmt.Errorf("error writing file to conditions file: %s", err.Error())
	}

	parameters := remoteConfig.Parameters
	parameterDir := filepath.Join(outputDir, config.ParametersDir)
	parameterFilePath := filepath.Join(parameterDir, config.ParametersFile)
	cs.fsClient.MkdirAll(parameterDir, 0744)

	parameterData, err := utils.JSONMarshal(parameters)
	if err != nil {
		return fmt.Errorf("error marshalling parameters: %s", err.Error())
	}

	parametersFile, err := cs.fsClient.OpenFile(parameterFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening parameter file for write: %s", err.Error())
	}
	_, err = parametersFile.Write(parameterData)
	if err != nil {
		return fmt.Errorf("error writing to parameter file: %s", err.Error())
	}
	return nil
}
func (cs *ClientStore) GetLocalConfig(dir string) (*remoteconfig.RemoteConfig, error) {
	remoteConfig := &remoteconfig.RemoteConfig{}
	conditionsFilePath := filepath.Join(dir, config.ConditionsDir)
	conditionsFromConfig, err := cs.readConditionsFromConfig(conditionsFilePath)
	if err != nil && err != io.EOF {
		return remoteConfig, err
	}
	remoteConfig.Conditions = conditionsFromConfig
	parametersDirPath := filepath.Join(dir, config.ParametersDir)
	parametersFromConfig, err := cs.readParametersFromConfig(parametersDirPath)
	if err != nil && err != io.EOF {
		return remoteConfig, err
	}
	remoteConfig.Parameters = parametersFromConfig
	return remoteConfig, nil
}
func (cs *ClientStore) ApplyConfig(rc remoteconfig.RemoteConfig, validateOnly bool) error {
	errs := utils.ValidateParameters(rc.Parameters)
	if len(errs) != 0 {
		errStringBuilder := strings.Builder{}
		for j := range errs {
			errStringBuilder.WriteString("\n\t" + errs[j].Error())
		}
		return fmt.Errorf("error validating parameter values. %s", errStringBuilder.String())
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
func (cs *ClientStore) GetRemoteConfigDiff(inputDir string) error {
	sourceConfig, err := cs.GetLocalConfig(inputDir)
	if err != nil {
		return err
	}
	remoteConfig, err := cs.GetLatestRemoteConfig()
	if err != nil {
		return err
	}
	utils.PrintDiff(*sourceConfig, *remoteConfig)
	return nil
}
func (cs *ClientStore) readConditionsFromConfig(dirPath string) ([]remoteconfig.Condition, error) {
	var conditions []remoteconfig.Condition
	conditionFilePath := filepath.Join(dirPath, config.ConditionsFile)
	file, err := cs.fsClient.OpenFile(conditionFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &conditions)
	if err != nil {
		return nil, err
	}
	return conditions, nil
}
func (cs *ClientStore) readParametersFromConfig(dirPath string) (map[string]remoteconfig.Parameter, error) {
	parameters := map[string]remoteconfig.Parameter{}
	s, err := cs.fsClient.Open(dirPath)
	if err != nil {
		return nil, err
	}
	files, err := s.Readdir(512)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}
		parameterFilePath := filepath.Join(dirPath, fileInfo.Name())
		file, err := cs.fsClient.OpenFile(parameterFilePath, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, &parameters)
		if err != nil {
			return nil, err
		}

	}

	return parameters, nil
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
		return nil, fmt.Errorf("Error while getting firebase app: %s", err)
	}

	return app, nil
}

func GetClientStore(ctx context.Context) (*ClientStore, error) {
	firebaseApp, err := getFirebaseApp(ctx)
	if err != nil {
		return nil, err
	}
	client, err := firebaseApp.RemoteConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting remoteconfig client: %s", err.Error())
	}
	return &ClientStore{remoteConfigClient: client, fsClient: afero.NewOsFs()}, nil
}
