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
		sParams[k] = v
	}
	for i := range remote.Parameters {
		rParams[i] = remote.Parameters[i]
	}
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("Generating diff for parameters")
	fmt.Println(GetRemoteDiffForParameters(sParams, rParams))
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")

}

const (
	Reset = "\u001B[0m"
	Red   = "\033[1;31m"
	Green = "\033[1;32m"
	Yellow= "\033[33m"
)

func GetRemoteDiffForConditions(source, remote []remoteconfig.Condition) string {
	diff := cmp.Diff(remote, source)
	greenDiff := strings.ReplaceAll(diff, "\n+", "\n"+Green)
	redDiff := strings.ReplaceAll(greenDiff, "\n-", "\n"+Red)
	finalDiff := strings.ReplaceAll(redDiff, "\n", Reset+"\n")
	return finalDiff
}

func GetRemoteDiffForParameters(source, remote map[string]remoteconfig.Parameter) string {
	diff := cmp.Diff(remote, source)
	maskedString := maskSecrets(diff)
	greenDiff := strings.ReplaceAll(maskedString, "\n+", "\n"+Green)
	redDiff := strings.ReplaceAll(greenDiff, "\n-", "\n"+Red)
	finalDiff := strings.ReplaceAll(redDiff, "\n", Reset+"\n")
	return finalDiff
}

func maskSecrets(input string) string {
    lines := strings.Split(input, "\n")
    maskedOutput := ""
    isMultiLineSecretJson := false

    for _, line := range lines {
        keyValue := strings.SplitN(line, ":", 2)
        if len(keyValue) > 1 {
            if isMultiLineSecretJson || strings.Contains(line, "SEC_") {
                keyValue[1] = "*******"
                if strings.Contains(line, "ExplicitValue"){
                    isMultiLineSecretJson = false
                }else{
                    isMultiLineSecretJson = true
                }
            }
        }
        maskedOutput += "\n" + strings.Join(keyValue, ":")
    }

    return maskedOutput
}
