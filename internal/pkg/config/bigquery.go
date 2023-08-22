package config

import (
	"google.golang.org/api/option"
)

// BigQuery ..
type BigQuery struct {
	DSN     string
	ADCPath string
}

// ToBigQueryClientOptions ..
func (s BigQuery) ToBigQueryClientOptions() []option.ClientOption {
	var opts []option.ClientOption
	if s.ADCPath != "" {
		opts = append(opts, option.WithCredentialsFile(s.ADCPath))
	}

	return opts
}

func bigQueryConfig() BigQuery {
	return BigQuery{
		DSN:     String("env.db.bigquery_dsn"),
		ADCPath: String("env.google_adc"),
	}
}
