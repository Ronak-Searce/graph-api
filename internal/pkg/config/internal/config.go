package internal

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	paramEnvKey = "APP_ENV"
	paramFmt    = "yaml"
	paramLocal  = "local.yaml"

	keyLogLevel    = "env.log_level"
	keyEnvID       = "app.env.id"
	keyServiceName = "env.service.name"
)

// LoadConfig ...
func LoadConfig(ctx context.Context) error {

	fPath := filepath.Join(".cfg", "k8s")
	fName := os.Getenv(paramEnvKey)
	if fName == "" {
		fName = paramLocal
	}

	viper.AddConfigPath(fPath)
	viper.SetConfigName(fName)
	viper.SetConfigType(paramFmt)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

// Values ...
func Values() map[string]any {
	out := map[string]any{}

	for _, k := range viper.AllKeys() {
		out[k] = viper.GetString(k)
	}

	return out
}

// Value ...
type Value struct {
	value any
}

// String ...
func (v Value) String() string {
	return cast.ToString(v.value)
}

// StringSlice ...
func (v Value) StringSlice() []string {
	return cast.ToStringSlice(v.value)
}

// IntSlice ...
func (v Value) IntSlice() []int {
	return cast.ToIntSlice(v.value)
}

// Bool ...
func (v Value) Bool() bool {
	return cast.ToBool(v.value)
}

// Int ...
func (v Value) Int() int {
	return cast.ToInt(v.value)
}

// Int32 ...
func (v Value) Int32() int32 {
	return cast.ToInt32(v.value)
}

// Uint32 ...
func (v Value) Uint32() uint32 {
	return cast.ToUint32(v.value)
}

// Int64 ...
func (v Value) Int64() int64 {
	return cast.ToInt64(v.value)
}

// Uint64 ...
func (v Value) Uint64() uint64 {
	return cast.ToUint64(v.value)
}

// Float64 ...
func (v Value) Float64() float64 {
	return cast.ToFloat64(v.value)
}

// Duration ...
func (v Value) Duration() time.Duration {
	return cast.ToDuration(v.value)
}

// ConfigValue ...
func ConfigValue(key string) Value {
	return Value{viper.Get(key)}
}

// IsConfigSet ...
func IsConfigSet(key string) bool {
	return viper.IsSet(key)
}

// LogLevel ...
func LogLevel() string {
	return viper.GetString(keyLogLevel)
}

// EnvID ...
func EnvID() string {
	return viper.GetString(keyEnvID)
}

// ServiceName ...
func ServiceName() string {
	return viper.GetString(keyServiceName)
}
