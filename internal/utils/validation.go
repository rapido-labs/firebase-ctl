package utils

import (
	"encoding/json"
	"fmt"
	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
	"strings"
)




func ValidateParameters(parameters map[string]remoteconfig.Parameter)[]error{
	errs := []error{}
	for k, v := range parameters{
		if v.Description ==""{
			errs = append(errs, fmt.Errorf("missing description for key: %s", k))
		}
		switch strings.ToLower(v.ValueType){
		case "string": continue
		case "json":
			err := validateJsonParameter(v)
			if err!= nil{
				errs = append(errs, fmt.Errorf("invalid json for key %s. error:%s",k, err.Error()))
			}
		default:
			errs = append(errs, fmt.Errorf("invalid value type for key:%s", k))
		}
	}
	return errs
}

func validateJsonParameter(parameter remoteconfig.Parameter)error{
	err := json.Unmarshal([]byte(parameter.DefaultValue.ExplicitValue), &map[string]interface{}{})
	if err!= nil{
		return fmt.Errorf("invalid json in default value. %s", err.Error())
	}
	for i,cv := range parameter.ConditionalValues{
		err := json.Unmarshal([]byte(cv.ExplicitValue), &map[string]interface{}{})
		if err!= nil{
			return fmt.Errorf("invalid json in conditional values. key:%s. error: %s", i, err.Error())
		}
	}
	return nil
}