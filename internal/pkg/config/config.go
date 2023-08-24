package config

import (
	"context"
)

// Config common configuration structure
type Config struct {
	AppEnv      string
	ServiceName string
	LogLevel    string
	EnvID       string
	Server      Server
	Spanner     Spanner
	BigQuery    BigQuery
	Pubsub      Pubsub
	AES         AES
}

// GetConfig initialise configuration and get Config object
func GetConfig(ctx context.Context) (Config, error) {
	err := Load(ctx)
	if err != nil {
		return Config{}, err
	}

	return Config{
		AppEnv:      String("app.env.id"),
		ServiceName: String("env.service.name"),
		LogLevel:    String("env.log_level"),
		EnvID:       String("app.env.id"),
		Server:      serverConfig(),
		Spanner:     spannerConfig(),
		BigQuery:    bigQueryConfig(),
		Pubsub:      pubsubConfig(),
		AES:         aesConfig(),
	}, nil
}

// IsProduction ...
func (c Config) IsProduction() bool {
	return c.AppEnv == "prod"
}
