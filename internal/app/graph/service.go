package graph

import (
	"context"

	graphInt "graph-api/internal/pkg/graph"

	loginPb "github.com/lyazii22/grpc-login/login/proto"
)

type ILoginProvider interface {
	CreateUser(ctx context.Context, in *loginPb.CreateUserRequest) (*loginPb.CreateUserResponse, error)
	GetUser(ctx context.Context, id string) (*loginPb.GetUserResponse, error)
	Login(ctx context.Context, in *loginPb.LoginRequest) (*loginPb.LoginResponse, error)
}

type Implementation struct {
	login ILoginProvider
}

// NewGraphAPI creates new graphql service instance
func NewGraphAPI(
	cliLogin ILoginProvider,
) *Implementation {
	return &Implementation{
		login: cliLogin,
	}
}

// Register registers resolvers implementation
func (i *Implementation) Register(reg *graphInt.Resolver) {
	// graphInt.InitCache(context.Background(), i)
	reg.InitImpl(i)
}
