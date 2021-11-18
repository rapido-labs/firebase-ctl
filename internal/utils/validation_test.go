package utils

import (
	"github.com/rapido-labs/firebase-ctl/internal/model"
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
	parameters := make(map[string]model.Parameter)
	parameters["validJson"] = model.Parameter{
		ConditionalValues: map[string]model.ParameterValue{"abcde": {
			ExplicitValue:   "{}",
			UseInAppDefault: false,
		}},
		DefaultValue: &model.ParameterValue{ExplicitValue: "{}"},
		Description:  "TestDescription",
		ValueType:    "json",
	}
	parameters["validArray"] = model.Parameter{
		ConditionalValues: map[string]model.ParameterValue{"abcde": {
			ExplicitValue:   "[]",
			UseInAppDefault: false,
		}},
		DefaultValue: &model.ParameterValue{ExplicitValue: "{}"},
		Description:  "TestDescription",
		ValueType:    "json",
	}
	parameters["validString"] = model.Parameter{
		ConditionalValues: nil,
		DefaultValue:      &model.ParameterValue{ExplicitValue: "adhfg"},
		Description:       "TestDescription",
		ValueType:         "string",
	}
	errs := ValidateParameters(parameters)
	assert.Len(c.T(), errs, 0)

	parameters["invalidJson"] = model.Parameter{
		ConditionalValues: nil,
		DefaultValue:      &model.ParameterValue{ExplicitValue: "adhfg"},
		Description:       "TestDesc",
		ValueType:         "json",
	}
	parameters["invalidType"] = model.Parameter{
		ConditionalValues: nil,
		DefaultValue:      &model.ParameterValue{ExplicitValue: "adhfg"},
		Description:       "TestDescription",
		ValueType:         "abc",
	}
	errs = ValidateParameters(parameters)
	assert.Len(c.T(), errs, 2)

	parameters["invalidJsonInConditionalValue"] = model.Parameter{
		ConditionalValues: map[string]model.ParameterValue{"abcde": model.ParameterValue{
			ExplicitValue:   "{",
			UseInAppDefault: false,
		}},
		DefaultValue: &model.ParameterValue{ExplicitValue: "{}"},
		Description:  "",
		ValueType:    "json",
	}
	errs = ValidateParameters(parameters)
	assert.Len(c.T(), errs, 3)

}
