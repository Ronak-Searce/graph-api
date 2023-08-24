# Graceful shutdown

[[_TOC_]]

## Goal

This library aims to make graceful shutdown easy and safe. It has some helpers to make setup easier.

## Helpers

### Graceful

- `GRPCServer` - pass `*grpc.Server` and it will stop it gracefully. New connections are not accepted, blocks until existing connections finish.
- `HTTPServer` - pass `*http.Server` and it will stop it gracefully. New connections are not accepted, blocks until existing connections finish.
- `PubSubTopic` - pass `*pubsub.Topic` and it will stop it gracefully. Publishes all remaining messages to publish.
- `Centrifuge` - pass `*centrifuge.Node` and it will stop it gracefully. New connections are not accepted, existing ones are disconnected with shutdown reason.
- `Redis` - pass `*redis.Client` and it will stop it gracefully.
- `Logger` - pass `*zap.SugaredLogger` and it will flush all buffered messages. Useful to ensure no logs are missing in logs explorer.
- `WaitGroup` - pass `*sync.WaitGroup` and it will wait for waitgroup to finish.
- `Tracer` - exports all remaining tracing spans.

### Graceless

- `Context` - pass `context.CancelFunc` and it will be called. Useful to cancel root context after all servers were stopped to ensure nothing is still invoking.
- `GRPCClient` - pass `*grpc.ClientConn` and it will stop it.
- `Spanner` - pass `*spanner.Client` and it will close all connections.

## Example

```go
package main

import (
	"context"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/go-redis/redis/v9"
	"gitlab.com/picnic-app/backend/libs/golang/graceful"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rc := redis.NewClient(&redis.Options{})
	rc.Get(ctx, "blabla")

	shutdown := graceful.New(
		&graceful.ShutdownManagerOptions{
			Timeout: time.Second * 15,
		},
		graceful.Parallel(
			&graceful.ParallelShutdownOptions{
				Name:    "servers",
				Timeout: time.Second * 30,
			},
		),
		graceful.Parallel(
			&graceful.ParallelShutdownOptions{
				Name:    "long tasks",
				Timeout: time.Second * 5,
			},
		),
		graceful.Context(cancel),
		graceful.Logger(nil),
		graceful.Parallel(
			&graceful.ParallelShutdownOptions{
				Name:    "clients",
				Timeout: time.Second * 5,
			},
			graceful.Redis(rc),
		),
	)
	shutdown.RegisterSignals(os.Interrupt, syscall.SIGTERM)

	time.Sleep(time.Second * 5)
	_ = shutdown.Shutdown(context.Background())

	for {
		runtime.Gosched()
	}
}
```

## Configure access to private modules on gitlab

This [notion page](https://www.notion.so/picnic-app/GitLab-How-to-use-GoLang-modules-acc4e52e21d34333a2c7e6f7ad263f33#2f798ec833984f3b84fada37b9563a11) 
explains how to configure local environment to work with private modules stored in Picnic gitlab repository.
