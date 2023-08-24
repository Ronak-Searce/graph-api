# healthcheck

[[_TOC_]]

## Goal

This library aims to make easy to set up healthchecks with prometheus metrics. There are some ready-to-use checks.

### Liveness vs Readiness

There are 2 probes in Kubernetes: liveness and readiness.

- Liveness - need to check whether pod is stuck and needs to be restarted. Most common way to do it - just respond `200 OK` on main port always without any subchecks. In most of our services GRPC is the main port so GRPC handler should be used. 
- Readiness - need to check whether pod is ready to accept traffic. There is no need to place in on main port so we use `debug` port for it.

## How to use

1. Create new `HealthChecker` with all required subchecks. You can also specify its name (for metrics) and enable metrics
1. Use it either as `http.Handler` (`HealthChecker` implements it) or call `GRPCRegisterService` method for GRPC handler

## How to add new checks

### In this library

1. Add new file with `Subcheck` implementation in `internal/pkg/healthcheck`
1. Add new file with `HealthCheckerOption` implementation in root directory

### Outside library

Implement `Subcheck` interface and add it via `AddSubcheck` method on `HealthChecker`.

```go
package main

import (
	"context"

	"gitlab.com/picnic-app/backend/libs/golang/healthcheck"
)

type BlablaSubcheck struct {}

func (*BlablaSubcheck) name() string { return "blabla" }
func (*BlablaSubcheck) run(_ context.Context) error { return nil }
func (*BlablaSubcheck) isWarning() bool { return false }

func main() {
	hc := healthcheck.New()
	hc.AddSubcheck(&BlablaSubcheck{})
}


```
