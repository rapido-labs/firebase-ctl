package firebase

import (
	"context"
	"encoding/json"
	"errors"
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
	latestRemoteConfig, err := cs.remoteConfigClient.GetRemoteConfig("")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return latestRemoteConfig.RemoteConfig, err
}
func (cs *ClientStore) BackupRemoteConfig(remoteConfig *remoteconfig.RemoteConfig, outputDir string) []error {
	errs := []error{}
	conditions := remoteConfig.Conditions
	conditionsDirPath := filepath.Join(outputDir, config.ConditionsDir)
	cs.fsClient.MkdirAll(conditionsDirPath, 0744)

	for _, condition := range conditions {
		data, err := utils.JSONMarshal(condition)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		fileName := filepath.Join(conditionsDirPath, condition.Name+".json")
		file, err := cs.fsClient.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
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
	cs.fsClient.MkdirAll(parametersDirPath, 0744)

	for key, parameter := range parameters {
		data, err := utils.JSONMarshal(parameter)
		if err != nil {
			errs = append(errs, err)
		}

		fileName := filepath.Join(parametersDirPath, key+".json")
		file, err := cs.fsClient.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
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
func (cs *ClientStore) GetLocalConfig(dir string) (*remoteconfig.RemoteConfig, error) {
	remoteConfig := &remoteconfig.RemoteConfig{}
	conditionsDirPath := filepath.Join(dir, config.ConditionsDir)
	conditionsFromConfig, err := cs.readConditionsFromConfig(conditionsDirPath)
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
func (cs *ClientStore) ApplyConfig(dir string, validateOnly bool) error {
	cfg, errs := cs.GetLocalConfig(dir)
	if errs != nil {
		return errs
	}
	template := remoteconfig.Template{
		Conditions:      cfg.Conditions,
		Parameters:      cfg.Parameters,
		ParameterGroups: cfg.ParameterGroups,
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
		return errors.New(fmt.Sprintf("error publishing template: %s ", err.Error()))
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
	s, err := cs.fsClient.Open(dirPath)
	if err != nil {
		return nil, err
	}
	files, err := s.Readdir(512)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			conditionFilePath := filepath.Join(dirPath, file.Name())
			condition := remoteconfig.Condition{}
			file, err := cs.fsClient.OpenFile(conditionFilePath, os.O_RDONLY, 0644)
			if err != nil {
				return nil, err
			}
			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bytes, &condition)
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, condition)
		}
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

	for _, file := range files {
		if !file.IsDir() {
			parameterFilePath := filepath.Join(dirPath, file.Name())
			parameter := remoteconfig.Parameter{}
			file, err := cs.fsClient.OpenFile(parameterFilePath, os.O_RDONLY, 0644)
			if err != nil {
				return nil, err
			}
			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(bytes, &parameter)
			if err != nil {
				return nil, err
			}
			fileBase := filepath.Base(file.Name())
			key := strings.TrimRight(fileBase, ".json")
			parameters[strings.Split(key, ".json")[0]] = parameter
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
