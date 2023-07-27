// package main

// import (
// "graph-api/internal/pkg/graph"
// graphPkg "graph-api/pkg/graph"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/99designs/gqlgen/graphql/handler"
// 	"github.com/99designs/gqlgen/graphql/playground"
// )

// const defaultPort = "8080"

// func main() {

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = defaultPort
// 	}

// 	srv := handler.NewDefaultServer(graphPkg.NewExecutableSchema(graphPkg.Config{Resolvers: &graph.Resolver{}}))

// 	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// 	http.Handle("/query", srv)

// 	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
// 	log.Fatal(http.ListenAndServe(":"+port, nil))
// }

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

	"gitlab.com/picnic-app/backend/libs/golang/config"
	"gitlab.com/picnic-app/backend/libs/golang/graceful"
	"gitlab.com/picnic-app/backend/libs/golang/healthcheck"
	"gitlab.com/picnic-app/backend/libs/golang/logger"
	"gitlab.com/picnic-app/backend/libs/golang/monitoring/monitoring"
	"gitlab.com/picnic-app/backend/libs/golang/monitoring/tracing"

	"graph-api/internal/app/graph"
	"graph-api/internal/core"

	// "graph-api/internal/pkg/auth/jwt"
	"graph-api/internal/pkg/client/login"
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

	loginCli := login.NewClient(ctx,
		viper.GetString("env.svc.login.host"),
		viper.GetInt("env.svc.login.port"),
		viper.GetBool("env.svc.login.secure"),
	)

	// jwtVerifier, err := jwt.NewVerifier(&jwt.Config{HMACSecretKey: viper.GetString("env.jwt.hmac_secret")})
	// if err != nil {
	// 	logger.Fatalf(ctx, "can't run app: %v", err)
	// }
	// authH := authPkg.NewHandler(jwtVerifier, authCli)

	graphAPI := graph.NewGraphAPI(
		loginCli,
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
