package main

import (
	"context"
	"log"

	"github.com/Ronak-Searce/graph-api/internal/pkg/config"
)

func main() {
	ctx := context.Background()
	err := config.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}
	println("env_id:", config.EnvID())
	println("service_name:", config.ServiceName())
	println("log_level:", config.LogLevel())
	println("env.test.str_value:", config.String("env.test.str_value"))
	println("env.test.int_value:", config.Int("env.test.int_value"))
	println("has value:", config.Has("env.test.int_value"))
	println("has no value:", config.Has("env.test.no_value"))
}
