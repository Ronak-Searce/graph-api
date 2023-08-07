package test

import (
	"fmt"
	gqlclient "graph-api/pkg/gql-client"
	"graph-api/pkg/graph"
	"testing"
)

func TestCreateUmpire(t *testing.T) {
	t.Log("check create group chat")

	// ctx := context.Background()

	// gql := gqlclient.NewGQLClient(&http.Client{
	// 	Timeout: 5 * time.Second,
	// })
	// anon := gql.Anonymous()
	username := getRandomString()
	fullname := getRandomString()
	email := getRandomEmail()
	password := getRandomString()

	body := fmt.Sprintf(`
		mutation {
			createUmpire(input:{username:"%s",password:"%s",email:"%s",name:"%s"}){
   				username
				name
				email
  			}
		}
	`, username, password, email, fullname)

	var resp struct {
		Data struct {
			CreateUmpire graph.Umpire `json:"createUmpire"`
		} `json:"data"`
	}

	err := gqlclient.Do(ctx, body, &resp)

	// users := helpCreateUsers(ctx, t, anon, nUsers)

	// user1, user2 := users[0], users[1]

	// chat := helpCreateGroupChat(ctx, t, "e2e-group-chat", user1, user2)
	// require.NotEmpty(t, chat.ID)
}
