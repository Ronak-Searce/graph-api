package config

import (
	"google.golang.org/api/option"
)

// Spanner ..
type Spanner struct {
	DSN     string
	ADCPath string
}

// ToSpannerClientOptions ..
// todo: checkup redundant dependency
func (s Spanner) ToSpannerClientOptions() []option.ClientOption {
	var opts []option.ClientOption
	if s.ADCPath != "" {
		opts = append(opts, option.WithCredentialsFile(s.ADCPath))
	}

	return opts
}

func spannerConfig() Spanner {
	return Spanner{
		DSN:     String("env.db.spanner_dsn"),
		ADCPath: String("env.google_adc"),
	}
}
