package graph

import (
	"context"

	graphInt "github.com/Ronak-Searce/graph-api/internal/pkg/graph"

	aflPb "github.com/Ronak-Searce/graph-api/api/proto/afl"
)

type IAflProvider interface {
	Login(ctx context.Context, in *aflPb.LoginRequest) (*aflPb.LoginResponse, error)
	CreateUmpire(ctx context.Context, in *aflPb.CreateUmpireRequest) (*aflPb.CreateUmpireResponse, error)
}

type Implementation struct {
	afl IAflProvider
}

// NewGraphAPI creates new graphql service instance
func NewGraphAPI(
	cliAfl IAflProvider,
) *Implementation {
	return &Implementation{
		afl: cliAfl,
	}
}

// Register registers resolvers implementation
func (i *Implementation) Register(reg *graphInt.Resolver) {
	// graphInt.InitCache(context.Background(), i)
	reg.InitImpl(i)
}
