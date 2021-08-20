package utils

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/rapido-labs/firebase-admin-go/v4/remoteconfig"
	"strings"
)

func PrintDiff(source, remote remoteconfig.RemoteConfig) {

	fmt.Println("Generating diff for conditions")
	fmt.Println(GetRemoteDiffForConditions(source.Conditions, remote.Conditions))
	sParams, rParams := map[string]remoteconfig.Parameter{}, map[string]remoteconfig.Parameter{}
	for k, v := range source.Parameters {
		v.ValueType=""
		sParams[k] = v
	}
	for i := range remote.Parameters {
		rParams[i] = remote.Parameters[i]
	}
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("Generating diff for parameters")
	fmt.Println(GetRemoteDiffForParameters(sParams, rParams))
}

const (
	Reset = "\u001B[0m"
	Red   = "\033[1;31m"
	Green = "\033[1;32m"
)

func GetRemoteDiffForConditions(source, remote []remoteconfig.Condition) string {
	sMap, rmap := map[string]remoteconfig.Condition{}, map[string]remoteconfig.Condition{}
	for _, v := range source {
		sMap[v.Name] = v
	}
	for _, v := range remote {
		rmap[v.Name] = v
	}
	diff := cmp.Diff(rmap, sMap)
	greenDiff := strings.ReplaceAll(diff, "\n+", "\n"+Green)
	redDiff := strings.ReplaceAll(greenDiff, "\n-", "\n"+Red)
	finalDiff := strings.ReplaceAll(redDiff, "\n", Reset+"\n")
	return finalDiff
}

func GetRemoteDiffForParameters(source, remote map[string]remoteconfig.Parameter) string {
	diff := cmp.Diff(remote, source)
	greenDiff := strings.ReplaceAll(diff, "\n+", "\n"+Green)
	redDiff := strings.ReplaceAll(greenDiff, "\n-", "\n"+Red)
	finalDiff := strings.ReplaceAll(redDiff, "\n", Reset+"\n")
	return finalDiff
}
