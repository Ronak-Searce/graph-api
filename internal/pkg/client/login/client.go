package login

import (
	"context"

	v1 "github.com/lyazii22/grpc-login/login/proto"
	"gitlab.com/picnic-app/backend/libs/golang/monitoring/tracing"

	grpcCli "graph-api/internal/pkg/client"
	// "graph-api/internal/pkg/model"
)

type Client struct {
	cli    v1.LoginClient
	tracer tracing.Tracer
}

func NewClient(ctx context.Context, host string, port int, secure bool) *Client {
	conn := grpcCli.NewClient(ctx, host, port, secure)
	return &Client{
		cli:    v1.NewLoginClient(conn),
		tracer: tracing.GetTracer("LoginClient"),
	}
}
