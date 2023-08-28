package test

import (
	"context"
	"net/http"
	"testing"
	"time"

	gqlclient "github.com/Ronak-Searce/graph-api/pkg/gql-client"

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
