package model

import (
	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
	"time"
)

// TagColor represents a tag color
type TagColor string

// Tag colors
const (
	colorUnspecified TagColor = ""
	Blue                      = "BLUE"
	Brown                     = "BROWN"
	Cyan                      = "CYAN"
	DeepOrange                = "DEEPORANGE"
	Green                     = "GREEN"
	Indigo                    = "INDIGO"
	Lime                      = "LIME"
	Orange                    = "ORANGE"
	Pink                      = "PINK"
	Purple                    = "PURPLE"
	Teal                      = "TEAL"
)

// Condition targets a specific group of users
// A list of these conditions make up part of a Remote Config template
type Condition struct {
	Expression string   `json:"expression"`
	Name       string   `json:"name"`
	TagColor   TagColor `json:"tagColor"`
}

// Config represents a Remote Config
type Config struct {
	Conditions      []Condition               `json:"conditions"`
	Parameters      map[string]Parameter      `json:"parameters"`
	ParameterGroups map[string]ParameterGroup `json:"parameterGroups"`
}

// Parameter .
type Parameter struct {
	ConditionalValues map[string]ParameterValue `json:"conditionalValues"`
	DefaultValue      *ParameterValue           `json:"defaultValue"`
	Description       string                    `json:"description"`
	ValueType         string                    `json:"valueType"`
}

// ParameterValue .
type ParameterValue struct {
	ExplicitValue   string `json:"value"`
	UseInAppDefault bool   `json:"useInAppDefault,omitempty"`
}

type ParameterGroup struct {
	Description string                `json:"description"`
	Parameters  map[string]*Parameter `json:"parameters"`
}

func ConvertToRemoteConfig(c Config) *remoteconfig.RemoteConfig {
	rc := &remoteconfig.RemoteConfig{
		Conditions: convertSourceConditionsToRemote(c.Conditions),
		Parameters: convertSourceParamsToRemote(c.Parameters),
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
		ParameterGroups: nil,
	}
	return rc
}

func ConvertToSourceConfig(c remoteconfig.RemoteConfig) *Config {
	rc := &Config{
		Conditions:      convertRemoteConditionsToSource(c.Conditions),
		Parameters:      convertRemoteParamsToSource(c.Parameters),
		ParameterGroups: nil,
	}
	return rc
}

func convertSourceConditionsToRemote(c []Condition) []remoteconfig.Condition {
	rcConditions := []remoteconfig.Condition{}

	for i := range c {
		rcConditions = append(rcConditions, remoteconfig.Condition{
			Expression: c[i].Expression,
			Name:       c[i].Name,
			TagColor:   remoteconfig.TagColor(c[i].TagColor),
		})
	}

	return rcConditions
}

func convertSourceParamsToRemote(p map[string]Parameter) map[string]remoteconfig.Parameter {
	rcParams := map[string]remoteconfig.Parameter{}

	for parameterKey, parameterValue := range p {
		cv := map[string]*remoteconfig.ParameterValue{}
		for conditionalValueKey, conditionalValueValue := range parameterValue.ConditionalValues {
			cv[conditionalValueKey] = &remoteconfig.ParameterValue{
				ExplicitValue:   conditionalValueValue.ExplicitValue,
				UseInAppDefault: conditionalValueValue.UseInAppDefault,
			}
		}
		if len(cv) == 0{
			cv = nil
		}
		rcParams[parameterKey] = remoteconfig.Parameter{
			ConditionalValues: cv,
			DefaultValue: &remoteconfig.ParameterValue{
				ExplicitValue:   p[parameterKey].DefaultValue.ExplicitValue,
				UseInAppDefault: p[parameterKey].DefaultValue.UseInAppDefault,
			},
			Description: p[parameterKey].Description,
		}
	}
	return rcParams
}

func convertRemoteConditionsToSource(c []remoteconfig.Condition) []Condition {
	var rcConditions []Condition

	for i := range c {
		rcConditions = append(rcConditions, Condition{
			Expression: c[i].Expression,
			Name:       c[i].Name,
			TagColor:   TagColor(c[i].TagColor),
		})
	}

	return rcConditions
}

func convertRemoteParamsToSource(p map[string]remoteconfig.Parameter) map[string]Parameter {
	rcParams := map[string]Parameter{}

	for oarameterKey, parameterValue := range p {
		cv := map[string]ParameterValue{}
		for conditionalValueKey, conditionalValueValue := range parameterValue.ConditionalValues {
			cv[conditionalValueKey] = ParameterValue{
				ExplicitValue:   conditionalValueValue.ExplicitValue,
				UseInAppDefault: conditionalValueValue.UseInAppDefault,
			}
		}
		if len(cv) == 0{
			cv = nil
		}
		rcParams[oarameterKey] = Parameter{
			ConditionalValues: cv,
			DefaultValue: &ParameterValue{
				ExplicitValue:   p[oarameterKey].DefaultValue.ExplicitValue,
				UseInAppDefault: p[oarameterKey].DefaultValue.UseInAppDefault,
			},
			Description: p[oarameterKey].Description,
		}
	}
	return rcParams
}
