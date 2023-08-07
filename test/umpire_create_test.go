package test

import (
	"context"
	gqlclient "graph-api/pkg/gql-client"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateUmpire(t *testing.T) {
	t.Log("check create group chat")

	ctx := context.Background()

	gql := gqlclient.NewGQLClient(&http.Client{
		Timeout: 5 * time.Second,
	})
	anon := gql.Anonymous()

	const nUsers = 2
	users := helpCreateUsers(ctx, t, anon, nUsers)

	user1, user2 := users[0], users[1]

	chat := helpCreateGroupChat(ctx, t, "e2e-group-chat", user1, user2)
	require.NotEmpty(t, chat.ID)
}
