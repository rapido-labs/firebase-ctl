package utils

import (
	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
}

func (c *ClientTestSuite)SetupTest(){}

func (c *ClientTestSuite)TestConditionsDiff()  {
	//nil and nil
	diff := GetRemoteDiffForConditions(nil, nil)
	assert.Equal(c.T(), "", diff)

	//for single identical-element array
	sourceArray := []remoteconfig.Condition{{
		Expression: "abcde",
		Name:       "name",
		TagColor:   "BLUE",
	}}
	diff = GetRemoteDiffForConditions(sourceArray, sourceArray)
	assert.Equal(c.T(), "", diff)

	//identical-but reversed array
	sourceArray = []remoteconfig.Condition{{
		Expression: "abcd",
		Name:       "name1",
		TagColor:   "BLUE",
	},{
		Expression: "efgh",
		Name:       "name2",
		TagColor:   "GREEN",
	}}
	remoteArray := []remoteconfig.Condition{{
		Expression: "efgh",
		Name:       "name2",
		TagColor:   "GREEN",
	},{
		Expression: "abcd",
		Name:       "name1",
		TagColor:   "BLUE",
	}}
	diff = GetRemoteDiffForConditions(sourceArray, remoteArray)
	assert.Equal(c.T(), "", diff)
}


func Test_Suite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

