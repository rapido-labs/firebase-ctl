package firebase

// import (
// 	"context"
// 	"os"
// 	"testing"

// 	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// type ClientMock struct {
// 	mock.Mock
// }

// func (c *ClientMock) RemoteConfig(ctx context.Context) (*remoteconfig.RemoteConfigClient, error) {
// 	args := c.Called(ctx)
// 	return args.Get(0).(*remoteconfig.RemoteConfigClient), args.Error(1)
// }

// type ClientTestSuite struct {
// 	suite.Suite
// 	client RemoteConfigClient
// }

// func (c *ClientTestSuite) Test_Client_wrong_credentials_file() {
// 	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/path/to/serviceaccount/json/file")
// 	ctx := context.Background()
// 	_, err := GetClientStore(ctx)

// 	assert.Contains(c.T(), err.Error(), "Error while getting remoteconfig client")
// }

// func (c *ClientTestSuite) Test_Client_return_remoteConfig_client() {
// 	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/path/to/serviceaccount/json/file")
// 	ctx := context.Background()
// 	remoteConfigClient := &remoteconfig.RemoteConfigClient{}
// 	c.clientMock.On("RemoteConfig", ctx).Return(&remoteConfigClient, nil)
// 	client, err := GetClientStore(ctx)

// 	assert.NoError(c.T(), err)
// 	assert.NotNil(c.T(), client)
// }

// func Test_Suite(t *testing.T) {
// 	suite.Run(t, new(ClientTestSuite))
// }
