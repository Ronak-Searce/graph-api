package core

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"

	"gitlab.com/picnic-app/backend/libs/golang/healthcheck"
	"gitlab.com/picnic-app/backend/libs/golang/monitoring/tracing"

	graphInt "github.com/Ronak-Searce/graph-api/internal/pkg/graph"

	"github.com/Ronak-Searce/graph-api/internal/app/graph"

	graphGen "github.com/Ronak-Searce/graph-api/pkg/graph"
)

// const banResponse = "You can not perform this operation"

// IGraphService ...
type IGraphService interface {
	Register(srv graphInt.IResolver)
}

type application struct {
	graphPort uint32
	graph     *handler.Server

	srv *http.Server
}

// New creates core application
func New(_ context.Context, httpPort uint32) *application {
	a := &application{
		graphPort: httpPort,
	}

	return a
}

// Run starts core application
func (a *application) Run(healthLive http.Handler, services ...*graph.Implementation) error {
	for _, service := range services {
		a.initGraphQL(service)
	}

	return a.listenGraphQL(a.graphPort, healthLive)
}

func (a *application) initGraphQL(impl *graph.Implementation) {
	reg := &graphInt.Resolver{}
	impl.Register(reg)

	srv := handler.New(graphGen.NewExecutableSchema(graphGen.Config{Resolvers: reg}))
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{
		MaxUploadSize: 1 << 30, // 1GiB
	})

	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	srv.Use(tracing.GraphqlInterceptor())
	// srv.Use(utils.GraphqlLoggerInterceptor())

	// srv.AroundResponses(a.mainGraphQLMiddleware)

	a.graph = srv
}

func (a *application) listenGraphQL(
	// authH *auth.Handler,
	port uint32,
	healthLive http.Handler,
) error {
	router := chi.NewRouter()
	router.Use(tracing.HTTPMiddleware)
	// router.Use(authH.HTTPMiddleware)

	router.Use(cors.New(cors.Options{
		AllowCredentials: true,
		Debug:            true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	}).Handler)

	var hh http.Handler = a.graph

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", hh)
	// router.Handle("/internal/query", authH.HTTPMetaMiddleware(hh))
	router.Handle(healthcheck.HTTPLivenessRoute, healthLive)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	a.srv = &http.Server{
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return a.srv.Serve(listener)
}

// mainGraphQLMiddleware checks for tokenID and isBanned fields
// func (a *application) mainGraphQLMiddleware(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
// 	if !graphql.HasOperationContext(ctx) {
// 		return next(ctx)
// 	}

// 	if op := graphql.GetOperationContext(ctx).Operation; op != nil && op.Operation == ast.Mutation {
// 		isBanned := utils.GetIsBannedFromCtx(ctx)
// 		if isBanned {
// 			return graphql.ErrorResponse(ctx, banResponse)
// 		}
// 	}

// 	return next(ctx)
// }

func (a *application) Close(ctx context.Context) error {
	if a.srv != nil {
		return a.srv.Shutdown(ctx)
	}

	return nil
}
