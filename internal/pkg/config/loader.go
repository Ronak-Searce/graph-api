package config

import (
	"context"
	"time"

	"github.com/Ronak-Searce/graph-api/internal/pkg/config/internal"
)

// Load initialise configuration
func Load(ctx context.Context) error {
	return internal.LoadConfig(ctx)
}

// Values returns a map of config key-value pairs
func Values() map[string]any {
	return internal.Values()
}

// Value get config value object
func Value(key string) internal.Value {
	return internal.ConfigValue(key)
}

// Get is simply a shorter version of Value()
func Get(key string) internal.Value {
	return Value(key)
}

// String ...
func String(key string) string {
	return internal.ConfigValue(key).String()
}

// StringSlice ...
func StringSlice(key string) []string {
	return internal.ConfigValue(key).StringSlice()
}

// IntSlice ...
func IntSlice(key string) []int {
	return internal.ConfigValue(key).IntSlice()
}

// Bool ...
func Bool(key string) bool {
	return internal.ConfigValue(key).Bool()
}

// Int ...
func Int(key string) int {
	return internal.ConfigValue(key).Int()
}

// Int32 ...
func Int32(key string) int32 {
	return internal.ConfigValue(key).Int32()
}

// Uint32 ...
func Uint32(key string) uint32 {
	return internal.ConfigValue(key).Uint32()
}

// Int64 ...
func Int64(key string) int64 {
	return internal.ConfigValue(key).Int64()
}

// Uint64 ...
func Uint64(key string) uint64 {
	return internal.ConfigValue(key).Uint64()
}

// Float64 ...
func Float64(key string) float64 {
	return internal.ConfigValue(key).Float64()
}

// Duration ...
func Duration(key string) time.Duration {
	return internal.ConfigValue(key).Duration()
}

// Has checks that key exists in config
func Has(key string) bool {
	return internal.IsConfigSet(key)
}

// LogLevel get configured log level
func LogLevel() string {
	return internal.LogLevel()
}

// EnvID env_id system value
func EnvID() string {
	return internal.EnvID()
}

// ServiceName service_name system value
func ServiceName() string {
	return internal.ServiceName()
}
