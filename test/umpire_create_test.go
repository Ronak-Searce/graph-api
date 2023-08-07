package test

import (
	"context"
	gqlclient "graph-api/pkg/gql-client"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUmpire(t *testing.T) {
	t.Log("check create umpire")

	ctx := context.Background()

	gql := gqlclient.NewGQLClient(&http.Client{
		Timeout: 5 * time.Second,
	})
	anon := gql.Anonymous()

	_, err := anon.CreateUser(ctx)

	assert.Nil(t, err, "fail to create umpire")
}
