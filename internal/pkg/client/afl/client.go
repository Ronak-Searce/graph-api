package afl

import (
	"context"

	aflPb "github.com/Ronak-Searce/graph-api/api/proto/afl"

	"github.com/Ronak-Searce/graph-api/internal/pkg/monitoring/tracing"

	grpcCli "github.com/Ronak-Searce/graph-api/internal/pkg/client"
	// "graph-api/internal/pkg/model"
)

type Client struct {
	cli    aflPb.AflClient
	tracer tracing.Tracer
}

func NewClient(ctx context.Context, host string, port int, secure bool) *Client {
	conn := grpcCli.NewClient(ctx, host, port, secure)
	return &Client{
		cli:    aflPb.NewAflClient(conn),
		tracer: tracing.GetTracer("LoginClient"),
	}
}

//	func (c *Client) CreateUser(ctx context.Context, in *loginPb.CreateUserRequest) (*loginPb.CreateUserResponse, error) {
//		res, err := c.cli.CreateUser(ctx, &loginPb.CreateUserRequest{
//			Username: in.Password,
//			Password: in.Password,
//		})
//		if err != nil {
//			return nil, err
//		}
//		return res, nil
//	}
//
//	func (c *Client) GetUser(ctx context.Context, id string) (*loginPb.GetUserResponse, error) {
//		res, err := c.cli.GetUser(ctx, &loginPb.GetUserRequest{Id: id}) // grpc client method
//		if err != nil {
//			return nil, err
//		}
//		return res, nil
//	}
func (c *Client) Login(ctx context.Context, in *aflPb.LoginRequest) (*aflPb.LoginResponse, error) {
	res, err := c.cli.Login(ctx, &aflPb.LoginRequest{
		Username: in.Username,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) CreateUmpire(ctx context.Context, in *aflPb.CreateUmpireRequest) (*aflPb.CreateUmpireResponse, error) {
	res, err := c.cli.CreateUmpire(ctx, &aflPb.CreateUmpireRequest{
		Username: in.Username,
		Fullname: in.Fullname,
		Password: in.Password,
		Email:    in.Email,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
