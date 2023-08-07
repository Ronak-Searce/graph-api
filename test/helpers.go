package test

import (
	"context"
	gqlclient "graph-api/pkg/gql-client"
	graph "graph-api/pkg/graph"
	"testing"

	"github.com/stretchr/testify/require"
)

func helpLogin(ctx context.Context, t *testing.T, u *gqlclient.User) graph.LoginOutput {
	t.Helper()
	resp, err := u.Login(ctx)
	require.NoError(t, err, "Login test fail")
	return graph.LoginOutput{
		Token:  resp.Data.Login.Token,
		Umpire: resp.Data.Login.Umpire,
	}
}
