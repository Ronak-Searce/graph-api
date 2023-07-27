package graph

import (
	"context"
	graphPkg "graph-api/pkg/graph"

	loginPb "github.com/lyazii22/grpc-login/login/proto"
)

func (i *Implementation) GetUser(ctx context.Context, id string) (*graphPkg.User, error) {
	res, err := i.login.GetUser(ctx, id) // grpc client method
	if err != nil {
		return nil, err
	}
	return &graphPkg.User{ID: res.Id, Name: res.Name}, nil
}
func (i *Implementation) CreateUser(ctx context.Context, input graphPkg.NewUser) (*graphPkg.User, error) {
	res, err := i.login.CreateUser(ctx, &loginPb.CreateUserRequest{
		Username: input.Password,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}
	return &graphPkg.User{ID: res.Id, Name: res.Name}, nil
}
func (i *Implementation) Login(ctx context.Context, input graphPkg.Login) (string, error) {
	res, err := i.login.Login(ctx, &loginPb.LoginRequest{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return "", err
	}
	return res.Token, nil
}

// imp <- Iresolver
// imp.login <- LoginClient
