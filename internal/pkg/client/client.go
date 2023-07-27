package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.com/picnic-app/backend/libs/golang/logger"
)

// common settings for connection
const (
	appName    = "graph-api"
	appVersion = "1.0"
)

type (
	ctxAppNameKey struct{}
)

func NewClient(ctx context.Context, host string, port int, secure bool) *grpc.ClientConn {
	return dial(ctx, host, port, secure, nil, nil)
}

func dial(
	ctx context.Context, host string, port int, secure bool,
	unaryInterceptors []grpc.UnaryClientInterceptor,
	streamInterceptors []grpc.StreamClientInterceptor,
) *grpc.ClientConn {
	cred := insecure.NewCredentials()

	if secure {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			logger.Fatal(ctx, err)
		}
		cred = credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
			RootCAs:    systemRoots,
		})
	}

	conn, err := grpc.DialContext(
		addAppNameToContext(ctx),
		fmt.Sprintf("%s:%d", host, port),
		[]grpc.DialOption{
			grpc.WithTransportCredentials(cred),
			// grpc.WithChainUnaryInterceptor(
			// 	append([]grpc.UnaryClientInterceptor{
			// 		NewUserIDUnaryInterceptor(),
			// 		tracing.ClientUnaryInterceptor(),
			// 	}, unaryInterceptors...)...,
			// ),
			// grpc.WithChainStreamInterceptor(
			// 	append([]grpc.StreamClientInterceptor{
			// 		NewUserIDStreamInterceptor(),
			// 		tracing.ClientStreamInterceptor(),
			// 	}, streamInterceptors...)...,
			// ),
			grpc.WithUserAgent(appName),
		}...,
	)
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}

	return conn
}

func addAppNameToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxAppNameKey{}, appName)
}
