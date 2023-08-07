package test

import(
	gapi "graph-api/pkg/gql-client"
	"context"
	"testing"
	"net/http"
	"time"
)


func TestCreateUmpire2(t *testing.T)
t.Log("check createUmpire")
t.skip("bug")

ctx := context.Background()

gql := gapi.NewGQLClient(&http.Client{
	Timeout: 5* time.Second,
})


