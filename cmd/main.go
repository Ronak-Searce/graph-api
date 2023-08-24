package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"github.com/Ronak-Searce/graph-api/internal/pkg/config"
	"github.com/Ronak-Searce/graph-api/internal/pkg/graceful"
	"github.com/Ronak-Searce/graph-api/internal/pkg/healthcheck"
	"github.com/Ronak-Searce/graph-api/internal/pkg/logger"
	"github.com/Ronak-Searce/graph-api/internal/pkg/monitoring/monitoring"
	"github.com/Ronak-Searce/graph-api/internal/pkg/monitoring/tracing"

	"github.com/Ronak-Searce/graph-api/internal/core"

	"github.com/Ronak-Searce/graph-api/internal/app/graph"

	// "graph-api/internal/pkg/auth/jwt"
	"github.com/Ronak-Searce/graph-api/internal/pkg/client/afl"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := config.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// TODO add healthchecker

	logger.InitWithString(config.LogLevel())

	// a := core.New(ctx)
	a := core.New(ctx, viper.GetUint32("service.ports.http"))

	// TODO: create interceptors for gql server
	// a.WithUnaryMW(mw.APIKeyAuthorizer)
	// a.WithUnaryMW(mw.LogInterceptor)

	// utils.InitURLCompleter(viper.GetString("env.svc.storage.url"))

	err = tracing.SetupExporter(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to set up tracing exporter: %v", err)
	}
	tracing.RegisterGoogleTracing()

	monitoring.RegisterPrometheusSuffix()

	aflCli := afl.NewClient(ctx,
		viper.GetString("env.svc.afl.host"),
		viper.GetInt("env.svc.afl.port"),
		viper.GetBool("env.svc.afl.secure"),
	)

	// jwtVerifier, err := jwt.NewVerifier(&jwt.Config{HMACSecretKey: viper.GetString("env.jwt.hmac_secret")})
	// if err != nil {
	// 	logger.Fatalf(ctx, "can't run app: %v", err)
	// }
	// authH := authPkg.NewHandler(jwtVerifier, authCli)

	graphAPI := graph.NewGraphAPI(
		aflCli,
	)

	healthReady := healthcheck.NewHealthCheck(
		healthcheck.WithName("ready"),
		healthcheck.WithMetrics(),
	)

	debugSrv, err := runDebugServer(ctx, healthReady)
	if err != nil {
		logger.Fatalf(ctx, "failed to start debug server: %v", err)
	}

	gracefulShutdown := graceful.New(
		&graceful.ShutdownManagerOptions{Timeout: 120 * time.Second},
		healthReady,
		graceful.Sleep(10*time.Second),
		graceful.Parallel(
			&graceful.ParallelShutdownOptions{
				Name:    "servers",
				Timeout: 120 * time.Second,
			},
			graceful.ShutdownContextErrorFunc(a.Close),
			graceful.HTTPServer(debugSrv),
		),
		graceful.Context(cancel),
		graceful.Parallel(
			&graceful.ParallelShutdownOptions{
				Name:    "clients",
				Timeout: 30 * time.Second,
			},
			graceful.Tracer(),
		),
		graceful.Logger(nil),
	)
	gracefulShutdown.RegisterSignals(os.Interrupt, syscall.SIGTERM)
	defer func() {
		_ = gracefulShutdown.Shutdown(context.Background())
	}()

	healthLive := healthcheck.NewHealthCheck(
		healthcheck.WithName("live"),
		healthcheck.WithMetrics(),
	)

	if err = a.Run(healthLive, graphAPI); err != nil {
		logger.Errorf(ctx, "can't run app: %v", err)
	}
}

func runDebugServer(ctx context.Context, healthReady http.Handler) (*http.Server, error) {
	port := config.String("service.ports.debug")
	addr := net.JoinHostPort("", port)

	listener, err := new(net.ListenConfig).Listen(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen debug port: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", monitoring.HTTPHandler())
	mux.Handle(healthcheck.HTTPReadinessRoute, healthReady)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		err := srv.Serve(listener)
		if err != nil {
			logger.Errorf(ctx, "debug server stopped: %v", err)
		}
	}()

	return srv, nil
}
