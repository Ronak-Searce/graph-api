package gqlclient

import (
	"context"
	"fmt"
)

// CreateUser ...
func (a *Anonymous) CreateUser(ctx context.Context) (*User, error) {

	u := NewUser(a.cli)

	if err := u.CreateUmpire(ctx); err != nil {
		return nil, fmt.Errorf("CreateUser: %v", err)
	}

	return u, nil
}
