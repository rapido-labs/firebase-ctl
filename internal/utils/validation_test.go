package utils

import (
	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ValidationTestSuite struct {
	suite.Suite
}

func (c *ValidationTestSuite) SetupTest() {}

func TestValidation(t *testing.T) {
	suite.Run(t, new(ValidationTestSuite))
}

func (c *ValidationTestSuite) TestParameters() {
	parameters := make(map[string]remoteconfig.Parameter)
	parameters["validJson"] = remoteconfig.Parameter{
		ConditionalValues: map[string]*remoteconfig.ParameterValue{"abcde": {
								ExplicitValue:   "{}",
								UseInAppDefault: false,
							}},
		DefaultValue: &remoteconfig.ParameterValue{ExplicitValue: "{}"},
		Description:  "TestDescription",
		ValueType:    "json",
	}
	parameters["validString"] = remoteconfig.Parameter{
		ConditionalValues: nil,
		DefaultValue:      &remoteconfig.ParameterValue{ExplicitValue: "adhfg"},
		Description:       "TestDescription",
		ValueType:         "string",
	}
	errs := ValidateParameters(parameters)
	assert.Len(c.T(), errs, 0)

	parameters["invalidJson"] = remoteconfig.Parameter{
		ConditionalValues: nil,
		DefaultValue:      &remoteconfig.ParameterValue{ExplicitValue: "adhfg"},
		Description:       "TestDesc",
		ValueType:         "json",
	}
	parameters["invalidType"] = remoteconfig.Parameter{
		ConditionalValues: nil,
		DefaultValue:      &remoteconfig.ParameterValue{ExplicitValue: "adhfg"},
		Description:       "TestDescription",
		ValueType:         "abc",
	}
	errs = ValidateParameters(parameters)
	assert.Len(c.T(), errs, 2)

	parameters["invalidJsonInConditionalValue"] = remoteconfig.Parameter{
		ConditionalValues: map[string]*remoteconfig.ParameterValue{"abcde": &remoteconfig.ParameterValue{
			ExplicitValue:   "{",
			UseInAppDefault: false,
		}},
		DefaultValue: &remoteconfig.ParameterValue{ExplicitValue: "{}"},
		Description:  "",
		ValueType:    "json",
	}
	errs = ValidateParameters(parameters)
	assert.Len(c.T(), errs, 3)

}
