package config

import (
	"fmt"
	"os"
)

type FirebaseConfig struct {
	Service_account_json_path string
}

func GetFirebaseConfig() (*FirebaseConfig, error) {
	service_account_json_path := os.Getenv(FIREBASE_AUTH_ENV_VAR)
	if service_account_json_path == "" {
		return nil, fmt.Errorf("%s env variable need to set", FIREBASE_AUTH_ENV_VAR)
	}
	return &FirebaseConfig{Service_account_json_path: service_account_json_path}, nil
}
