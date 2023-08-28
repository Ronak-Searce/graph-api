package gqlclient

import (
	"context"
	"fmt"
	"graph-api/pkg/graph"
)

func (u *User) CreateUmpire(ctx context.Context) error {

	body := fmt.Sprintf(`
		mutation {
			createUmpire(input:{username:"%s",password:"%s",email:"%s",name:"%s"}){
   				username
				name
				email
  			}
		}
	`, u.Username, u.Password, u.Email, u.Fullname)

	var resp struct {
		Data struct {
			CreateUmpire graph.Umpire `json:"createUmpire"`
		} `json:"data"`
	}

	err := u.cli.Do(ctx, body, &resp)

	if err != nil {
		return err
	}

	return nil
}
