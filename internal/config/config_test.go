package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (c *ConfigTestSuite) SetupTest() {
	clearenv()
}

func (c *ConfigTestSuite) Test_Config_env_var_not_present() {
	_, err := GetFirebaseConfig()
	assert.Error(c.T(), err, "GOOGLE_APPLICATION_CREDENTIALS env var need to set")
}

func (c *ConfigTestSuite) Test_Config_env_presnt() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "path/to/json/file")
	_, err := GetFirebaseConfig()
	assert.NoError(c.T(), err)
}

func Test_Suite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func clearenv() {
	os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
}
