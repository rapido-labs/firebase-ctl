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
	return args.Get(0).(*remoteconfig.Response), args.Error(1)
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
	rc := ClientStore{RemoteConfigClient: c.mock}
	c.mock.On("GetRemoteConfig",mock.Anything).Return(&remoteconfig.Response{
		RemoteConfig: remoteconfig.RemoteConfig{},
		Etag:         "",
	}, nil).Times(1)
	cfg, err:= rc.GetLatestRemoteConfig(context.Background())
	assert.NoError(c.T(), err)
	assert.NotNil(c.T(), cfg)
	c.mock.AssertExpectations(c.T())
}
func (c *ClientTestSuite) TestGetRemoteConfigClientErrorsOut() {
	rc := ClientStore{RemoteConfigClient: c.mock}
	c.mock.On("GetRemoteConfig",mock.Anything).Return((*remoteconfig.Response)(nil) , errors.New("test error")).Times(1)
	cfg, err:= rc.GetLatestRemoteConfig(context.Background())
	assert.Error(c.T(), err)
	assert.Equal(c.T(), err.Error(), "test error")
	assert.Nil(c.T(), cfg)
	c.mock.AssertExpectations(c.T())


}

func (c *ClientTestSuite)TestBackupToFsSuccess()  {
	tempFs := afero.NewMemMapFs()
	cs := ClientStore{RemoteConfigClient: c.mock, FsClient: tempFs}
	dummyResponse := &remoteconfig.Response{
		RemoteConfig: remoteconfig.RemoteConfig{
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
	errs:= cs.BackupRemoteConfig(dummyResponse,outputDir)
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
	cs := ClientStore{RemoteConfigClient: c.mock, FsClient: tempFs}
	dummyResponse := &remoteconfig.Response{
		RemoteConfig: remoteconfig.RemoteConfig{
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


	cs.FsClient =tempFs

	errs := cs.BackupRemoteConfig(dummyResponse, outputDir)
	assert.Len(c.T(), errs, 2)
}





func (c *ClientTestSuite) Test_Client_return_remoteConfig_client() {
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/path/to/serviceaccount/json/file")
	//ctx := context.Background()
	//remoteConfigClient := &remoteconfig.RemoteConfigClient{}
	//c.clientMock.On("RemoteConfig", ctx).Return(&remoteConfigClient, nil)
	//client, err := GetClientStore(ctx)
	//
	//assert.NoError(c.T(), err)
	//assert.NotNil(c.T(), client)
}

func Test_Suite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}


