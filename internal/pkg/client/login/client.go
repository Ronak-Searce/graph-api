package login

import (
	"context"

	loginPb "github.com/lyazii22/grpc-login/login/proto"
	"gitlab.com/picnic-app/backend/libs/golang/monitoring/tracing"

	grpcCli "graph-api/internal/pkg/client"
	// "graph-api/internal/pkg/model"
)

type Client struct {
	cli    loginPb.LoginClient
	tracer tracing.Tracer
}

func NewClient(ctx context.Context, host string, port int, secure bool) *Client {
	conn := grpcCli.NewClient(ctx, host, port, secure)
	return &Client{
		cli:    loginPb.NewLoginClient(conn),
		tracer: tracing.GetTracer("LoginClient"),
	}
}

func (c *Client) CreateUser(ctx context.Context, in *loginPb.CreateUserRequest) (*loginPb.CreateUserResponse, error) {
	res, err := c.cli.CreateUser(ctx, &loginPb.CreateUserRequest{
		Username: in.Password,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (c *Client) GetUser(ctx context.Context, id string) (*loginPb.GetUserResponse, error) {
	res, err := c.cli.GetUser(ctx, &loginPb.GetUserRequest{Id: id}) // grpc client method
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (c *Client) Login(ctx context.Context, in *loginPb.LoginRequest) (*loginPb.LoginResponse, error) {
	res, err := c.cli.Login(ctx, &loginPb.LoginRequest{
		Username: in.Username,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// func (i *Implementation) GetUser(ctx context.Context, id string) (*graphPkg.User, error) {
// 	res, err := i.login.GetUser(ctx, id) // grpc client method
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &graphPkg.User{ID: res.Id, Name: res.Name}, nil
// }
// func (c *Client) CreateUser(ctx context.Context, input graphPkg.NewUser) (*graphPkg.User, error) {
// 	res, err := i.login.CreateUser(ctx, &loginPb.CreateUserRequest{
// 		Username: input.Password,
// 		Password: input.Password,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &graphPkg.User{ID: res.Id, Name: res.Name}, nil
// }
// func (i *Implementation) Login(ctx context.Context, input graphPkg.Login) (string, error) {
// 	res, err := i.login.Login(ctx, &loginPb.LoginRequest{
// 		Username: input.Username,
// 		Password: input.Password,
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	return res.Token, nil
// }
