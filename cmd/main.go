// package main

// import (
	// "graph-api/internal/pkg/graph"
	graphPkg "graph-api/pkg/graph"
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
	authV1 "gitlab.com/picnic-app/backend/libs/golang/protobuf-registry/gen/auth-api/auth/v1"

	"graph-api/internal/app/graph"
	"graph-api/internal/core"
	authPkg "graph-api/internal/pkg/auth"
	"graph-api/internal/pkg/auth/jwt"
	grpcCli "graph-api/internal/pkg/client"
	"graph-api/internal/pkg/client/admin"
	"graph-api/internal/pkg/client/analytics"
	"graph-api/internal/pkg/client/app"
	"graph-api/internal/pkg/client/blacklist"
	"graph-api/internal/pkg/client/chat"
	"graph-api/internal/pkg/client/circle"
	"graph-api/internal/pkg/client/content"
	"graph-api/internal/pkg/client/document"
	"graph-api/internal/pkg/client/election"
	"graph-api/internal/pkg/client/feed"
	"graph-api/internal/pkg/client/growth"
	"graph-api/internal/pkg/client/languages"
	"graph-api/internal/pkg/client/linkmetadata"
	"graph-api/internal/pkg/client/messaging"
	"graph-api/internal/pkg/client/notification"
	"graph-api/internal/pkg/client/profile"
	"graph-api/internal/pkg/client/recommendation"
	"graph-api/internal/pkg/client/report"
	"graph-api/internal/pkg/client/seeds"
	"graph-api/internal/pkg/client/slice"
	"graph-api/internal/pkg/client/storage"
	"graph-api/internal/pkg/utils"
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

	utils.InitURLCompleter(viper.GetString("env.svc.storage.url"))

	err = tracing.SetupExporter(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to set up tracing exporter: %v", err)
	}
	tracing.RegisterGoogleTracing()

	monitoring.RegisterPrometheusSuffix()

	authCli := authV1.NewAuthAPIClient(
		grpcCli.NewClient(
			ctx,
			viper.GetString("env.svc.auth.host"),
			viper.GetInt("env.svc.auth.port"),
			viper.GetBool("env.svc.auth.secure"),
		),
	)
	appCli := app.NewClient(ctx,
		viper.GetString("env.svc.app.host"),
		viper.GetInt("env.svc.app.port"),
		viper.GetBool("env.svc.app.secure"),
	)
	circleCli := circle.NewClient(ctx, viper.GetString("env.svc.circle.host"),
		viper.GetInt("env.svc.circle.port"),
		viper.GetBool("env.svc.circle.secure"))
	contentCli := content.NewClient(ctx, viper.GetString("env.svc.content.host"),
		viper.GetInt("env.svc.content.port"),
		viper.GetBool("env.svc.content.secure"))
	profileCli := profile.NewClient(ctx, viper.GetString("env.svc.profile.host"),
		viper.GetInt("env.svc.profile.port"),
		viper.GetBool("env.svc.profile.secure"))
	storageCli := storage.NewClient(ctx, viper.GetString("env.svc.storage.host"),
		viper.GetInt("env.svc.storage.port"),
		viper.GetBool("env.svc.storage.secure"),
	)
	chatCli := chat.NewClient(ctx, viper.GetString("env.svc.chat.host"),
		viper.GetInt("env.svc.chat.port"),
		viper.GetBool("env.svc.chat.secure"))
	seedsCli := seeds.NewClient(ctx, viper.GetString("env.svc.seeds.host"),
		viper.GetInt("env.svc.seeds.port"),
		viper.GetBool("env.svc.seeds.secure"))
	feedCli := feed.NewClient(ctx, viper.GetString("env.svc.feed.host"),
		viper.GetInt("env.svc.feed.port"),
		viper.GetBool("env.svc.feed.secure"))
	electionCli := election.NewClient(ctx, viper.GetString("env.svc.election.host"),
		viper.GetInt("env.svc.election.port"),
		viper.GetBool("env.svc.election.secure"))
	languagesCli := languages.NewClient(ctx, viper.GetString("env.svc.languages.host"),
		viper.GetInt("env.svc.languages.port"),
		viper.GetBool("env.svc.languages.secure"))
	blacklistCli := blacklist.NewClient(ctx, viper.GetString("env.svc.blacklist.host"),
		viper.GetInt("env.svc.blacklist.port"),
		viper.GetBool("env.svc.blacklist.secure"))
	reportCli := report.NewClient(ctx, viper.GetString("env.svc.report.host"),
		viper.GetInt("env.svc.report.port"),
		viper.GetBool("env.svc.report.secure"))
	documentCli := document.NewClient(ctx, viper.GetString("env.svc.document.host"),
		viper.GetInt("env.svc.document.port"),
		viper.GetBool("env.svc.document.secure"))
	notificationCli := notification.NewClient(ctx, viper.GetString("env.svc.notification.host"),
		viper.GetInt("env.svc.notification.port"),
		viper.GetBool("env.svc.notification.secure"))
	adminCli := admin.NewClient(ctx, viper.GetString("env.svc.admin.host"),
		viper.GetInt("env.svc.admin.port"),
		viper.GetBool("env.svc.admin.secure"))
	sliceCli := slice.NewClient(ctx, viper.GetString("env.svc.slice.host"),
		viper.GetInt("env.svc.slice.port"),
		viper.GetBool("env.svc.slice.secure"))
	messagingCli := messaging.NewClient(ctx, viper.GetString("env.svc.messaging.host"),
		viper.GetInt("env.svc.messaging.port"),
		viper.GetBool("env.svc.messaging.secure"))
	analyticsCli := analytics.NewClient(ctx, viper.GetString("env.svc.analytics.host"),
		viper.GetInt("env.svc.analytics.port"),
		viper.GetBool("env.svc.analytics.secure"))
	linkmetadataCli := linkmetadata.NewClient(ctx, viper.GetString("env.svc.linkmetadata.host"),
		viper.GetInt("env.svc.linkmetadata.port"),
		viper.GetBool("env.svc.linkmetadata.secure"))
	recommendationCli := recommendation.NewClient(ctx,
		viper.GetString("env.svc.recommendation.host"),
		viper.GetInt("env.svc.recommendation.port"),
		viper.GetBool("env.svc.recommendation.secure"),
	)
	growthCli := growth.NewClient(ctx,
		viper.GetString("env.svc.growth.host"),
		viper.GetInt("env.svc.growth.port"),
		viper.GetBool("env.svc.growth.secure"),
	)

	jwtVerifier, err := jwt.NewVerifier(&jwt.Config{HMACSecretKey: viper.GetString("env.jwt.hmac_secret")})
	if err != nil {
		logger.Fatalf(ctx, "can't run app: %v", err)
	}
	authH := authPkg.NewHandler(jwtVerifier, authCli)

	graphAPI := graph.NewGraphAPI(
		authCli,
		circleCli,
		contentCli,
		profileCli,
		storageCli,
		chatCli,
		feedCli,
		electionCli,
		seedsCli,
		languagesCli,
		blacklistCli,
		reportCli,
		documentCli,
		notificationCli,
		adminCli,
		sliceCli,
		messagingCli,
		analyticsCli,
		appCli,
		linkmetadataCli,
		recommendationCli,
		growthCli,
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

	if err = a.Run(authH, healthLive, graphAPI); err != nil {
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
