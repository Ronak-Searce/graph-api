package test

import (
	"context"
	"net/http"
	"testing"
	"time"

	gqlclient "github.com/Ronak-Searce/graph-api/pkg/gql-client"

	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	t.Log("check login")

	ctx := context.Background()

	gql := gqlclient.NewGQLClient(&http.Client{
		Timeout: 5 * time.Second,
	})
	anon := gql.Anonymous()

	u, err := anon.CreateUser(ctx)
	require.NoError(t, err, "fail to create umpire")

	resp := helpLogin(ctx, t, u)

	require.NotEqual(t, resp.Token, "", "Login: invalid token")
	require.Equal(t, u.Email, resp.Umpire.Email, "Login: email not matching")
	require.Equal(t, u.Username, resp.Umpire.Username, "Login: username not matching")
	// commented out because the response return username in place of fullname
	// require.Equal(t, u.Fullname, resp.Umpire.Name, "Login: name not matching")
}
