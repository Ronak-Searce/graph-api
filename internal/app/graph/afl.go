package graph

import (
	"context"
	graphPkg "graph-api/pkg/graph"

	aflPb "graph-api/api/proto/afl"
)

func (i *Implementation) Login(ctx context.Context, input graphPkg.LoginInput) (*graphPkg.LoginOutput, error) {
	res, err := i.afl.Login(ctx, &aflPb.LoginRequest{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}
	return &graphPkg.LoginOutput{Token: res.AccessToken, Umpire: &graphPkg.Umpire{
		Name:     res.Umpire.Username,
		Email:    res.Umpire.Email,
		Username: res.Umpire.Username,
	}}, nil
}

func (i *Implementation) CreateUmpire(ctx context.Context, input graphPkg.NewUmpire) (*graphPkg.Umpire, error) {
	res, err := i.afl.CreateUmpire(ctx, &aflPb.CreateUmpireRequest{
		Username: input.Username,
		Fullname: input.Name,
		Email:    input.Email,
	})
	if err != nil {
		return nil, err
	}
	return &graphPkg.Umpire{Name: res.Umpire.Fullname,
		Username: res.Umpire.Username,
		Email:    res.Umpire.Email,
	}, nil
}
