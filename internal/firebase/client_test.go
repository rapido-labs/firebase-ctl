package firebase

import (
	"context"
	"errors"
	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) GetRemoteConfig(versionNumber string) (*remoteconfig.Response, error) {
	args := c.Called(versionNumber)
	if args.Get(0)==nil{
		return nil, args.Error(1)
	}
	return args.Get(0).(*remoteconfig.Response), args.Error(1)
}
func (c * ClientMock)PublishTemplate(ctx context.Context, template remoteconfig.Template, validateOnly bool)(*remoteconfig.Template,error){
	args := c.Called(ctx, template, validateOnly)
	if args.Get(0)==nil{
		return nil, args.Error(1)
	}
	return args.Get(0).(*remoteconfig.Template), args.Error(1)
}

func (c *ClientMock) RemoteConfig(ctx context.Context) (*remoteconfig.Client, error) {
	args := c.Called(ctx)
	return args.Get(0).(*remoteconfig.Client), args.Error(1)
}

type ClientTestSuite struct {
	suite.Suite
	mock *ClientMock
}

func (c *ClientTestSuite)SetupTest(){
	c.mock =new(ClientMock)
}

func (c *ClientTestSuite) TestGetRemoteConfigClientReturnsValidConfig() {
	rc := ClientStore{remoteConfigClient: c.mock}
	c.mock.On("GetRemoteConfig",mock.Anything).Return(&remoteconfig.Response{
		RemoteConfig: &remoteconfig.RemoteConfig{
			Conditions:      nil,
			Parameters:      nil,
			Version:         remoteconfig.Version{},
			ParameterGroups: nil,
		},
		Etag:         "",
	}, nil).Times(1)
	cfg, err:= rc.GetLatestRemoteConfig()
	assert.NoError(c.T(), err)
	assert.NotNil(c.T(), cfg)
	c.mock.AssertExpectations(c.T())
}
func (c *ClientTestSuite) TestGetRemoteConfigClientErrorsOut() {
	rc := ClientStore{remoteConfigClient: c.mock}
	c.mock.On("GetRemoteConfig",mock.Anything).Return((*remoteconfig.Response)(nil) , errors.New("test error")).Times(1)
	cfg, err:= rc.GetLatestRemoteConfig()
	assert.Error(c.T(), err)
	assert.Equal(c.T(), err.Error(), "test error")
	assert.Nil(c.T(), cfg)
	c.mock.AssertExpectations(c.T())


}

func (c *ClientTestSuite)TestBackupToFsSuccess()  {
	tempFs := afero.NewMemMapFs()
	cs := ClientStore{remoteConfigClient: c.mock, fsClient: tempFs}
	dummyResponse := &remoteconfig.Response{
		RemoteConfig: &remoteconfig.RemoteConfig{
			Conditions:      []remoteconfig.Condition{remoteconfig.Condition{
				Expression: "hello",
				Name:       "there",
				TagColor:   remoteconfig.Blue,
			}},
			Parameters:      map[string]remoteconfig.Parameter{"hello": {
				ConditionalValues: nil,
				DefaultValue:      nil,
				Description:       "",
						}},
			Version:         remoteconfig.Version{},
			ParameterGroups: nil,
		},
		Etag:         "",
	}
	outputDir := "sample/outputdir"
	errs:= cs.BackupRemoteConfig(dummyResponse.RemoteConfig,outputDir)
	assert.Lenf(c.T(), errs, 0, "error length does not match. expected 0")
	file, err := tempFs.OpenFile(filepath.Join(outputDir, "conditions", "there.json"),os.O_RDONLY, 0644)
	if err!= nil{
		c.T().Fail()
	}

	contents, err := ioutil.ReadAll(file)
	if err!=nil{
		c.T().Fail()
	}
	assert.Contains(c.T(), string(contents), "there")
	assert.Contains(c.T(), string(contents), "hello")
	assert.Contains(c.T(), string(contents), remoteconfig.Blue)
}


func (c *ClientTestSuite)TestBackupToFsFailureBecauseFSErrors()  {
	memFs := afero.NewMemMapFs()
	tempFs := afero.NewReadOnlyFs(memFs)
	cs := ClientStore{remoteConfigClient: c.mock, fsClient: tempFs}
	dummyResponse := &remoteconfig.Response{
		RemoteConfig: &remoteconfig.RemoteConfig{
			Conditions:      []remoteconfig.Condition{remoteconfig.Condition{
				Expression: "hello",
				Name:       "there",
				TagColor:   remoteconfig.Blue,
			}},
			Parameters:      map[string]remoteconfig.Parameter{"hello": {
				ConditionalValues: nil,
				DefaultValue:      nil,
				Description:       "",
			}},
			Version:         remoteconfig.Version{},
			ParameterGroups: nil,
		},
		Etag:         "",
	}
	outputDir := "sample/outputdir"


	cs.fsClient =tempFs

	errs := cs.BackupRemoteConfig(dummyResponse.RemoteConfig, outputDir)
	assert.Len(c.T(), errs, 2)
}

func (c * ClientTestSuite)TestApplyConfig(){
	tempFs := afero.NewOsFs()
	cs := ClientStore{fsClient: tempFs, remoteConfigClient:c.mock}
	//successful publish
	c.mock.On("PublishTemplate", context.Background(),mock.Anything, false).Return(&remoteconfig.Template{}, nil).Times(1)
	err := cs.ApplyConfig("./test", false)
	assert.NoError(c.T(), err)
	c.mock.AssertExpectations(c.T())

	// invalid directory given
	err = cs.ApplyConfig("./test1", false)

	assert.Equal(c.T(), "open test1/conditions: no such file or directory", err.Error())
	c.mock.AssertExpectations(c.T())

	c.mock.On("PublishTemplate", context.Background(),mock.Anything, false).Return((*remoteconfig.Template)(nil), errors.New("test error")).Times(1)
	err = cs.ApplyConfig("./test", false)
	assert.Contains(c.T(), err.Error(),"test error")
	c.mock.AssertExpectations(c.T())


}

func (c *ClientTestSuite)TestGetDiff(){
	tempFs := afero.NewOsFs()
	cs := ClientStore{fsClient: tempFs, remoteConfigClient:c.mock}
	c.mock.On("GetRemoteConfig", "").Return(&remoteconfig.Response{RemoteConfig:&remoteconfig.RemoteConfig{
		Conditions:      []remoteconfig.Condition{},
		Parameters:      make(map[string]remoteconfig.Parameter),
		Version:         remoteconfig.Version{},
		ParameterGroups: nil,
	}}, nil).Times(1)

	// successfully find the diff
	err := cs.GetRemoteConfigDiff("./test")
	assert.NoError(c.T(), err)
	c.mock.AssertExpectations(c.T())

	//pass an invalid directory
	err = cs.GetRemoteConfigDiff("./test1")
	assert.Contains(c.T(), err.Error(), "no such file or directory")
	c.mock.AssertExpectations(c.T())

	//google api returns an error
	c.mock.On("GetRemoteConfig", "").Return(nil, errors.New("test error")).Times(1)
	err = cs.GetRemoteConfigDiff("./test")
	assert.Contains(c.T(), err.Error(), "test error")
	// successfully find the diff
	c.mock.On("GetRemoteConfig", "").Return(nil, errors.New("test error")).Times(1)
	err = cs.GetRemoteConfigDiff("./test")
	assert.Contains(c.T(), err.Error(), "test error")
	c.mock.AssertExpectations(c.T())

}
func (c *ClientTestSuite)TestGetLocalConfig(){
	cs := &ClientStore{fsClient: afero.NewOsFs()}
	rc, err := cs.GetLocalConfig("./test")
	assert.Nil(c.T(), err, "error was not expected")
	assert.Len(c.T(), rc.Conditions, 2, "unexpected conditions length")
	assert.Len(c.T(), rc.Parameters, 4, "unexpected parameters length")

}

func (c *ClientTestSuite)TestBackup(){
	configToWrite := remoteconfig.RemoteConfig{
		Conditions:      []remoteconfig.Condition{remoteconfig.Condition{
			Expression: "a==b",
			Name:       "test_name",
			TagColor:   remoteconfig.Blue,
		}},
		Parameters:      map[string]remoteconfig.Parameter{
			"test":remoteconfig.Parameter{
				ConditionalValues: map[string]*remoteconfig.ParameterValue{},
				DefaultValue:      &remoteconfig.ParameterValue{
					ExplicitValue:   "test_value",
					UseInAppDefault: true,
				},
				Description:       "test_description",
			},
		},
		Version:         remoteconfig.Version{},
		ParameterGroups: nil,
	}
	cs := ClientStore{
		remoteConfigClient: c.mock,
		fsClient:           afero.NewMemMapFs(),
	}
	errs := cs.BackupRemoteConfig(&configToWrite,"test")
	assert.Len(c.T(), errs, 0)

	localConfig, err := cs.GetLocalConfig("test")
	assert.NoError(c.T(), err)
	assert.Equal(c.T(),configToWrite, *localConfig)

}


func Test_Suite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}


