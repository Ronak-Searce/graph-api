package gqlclient

import (
	"context"
	"fmt"

	graph "github.com/Ronak-Searce/graph-api/pkg/graph"
)

type LoginResponse struct {
	Data struct {
		Login struct {
			Token  string        `json:"token"`
			Umpire *graph.Umpire `json:"umpire"`
		} `json:"login"`
	} `json:"data"`
}

func (u *User) Login(ctx context.Context) (*LoginResponse, error) {

	body := fmt.Sprintf(`
		mutation {
			login(input:{username:"%s",password:"%s"}){
    			token
   				umpire{
     			 name
     			 email
     			 username
    			}	
  			}
		}
	`, u.Username, u.Password)
	var resp LoginResponse

	if err := u.cli.Do(ctx, body, &resp); err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	return &resp, nil
}
