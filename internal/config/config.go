package config

import (
	"os"

	"graph-api/internal/pkg/errors"
)

type configParam = string

// ConfJWTHMACSecret environment variable names
const (
	ConfJWTHMACSecret configParam = "JWT_HMAC_SECRET" //nolint:gosec
)

// GetEnv get environment value by key, returns error if key not found
func GetEnv(key configParam) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", errors.Errorf(errors.UnknownEnvKey, key)
	}
	return val, nil
}
