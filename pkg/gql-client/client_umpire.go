package gqlclient

import (
	"fmt"
	"time"
)

type Umpire struct {
	cli      *GQLClient
	Fullname string
	Username string
	Email    string
}

func NewUmpire(cli *GQLClient) *Umpire {
	value := time.Now().UTC().UnixMicro()
	fullName := fmt.Sprintf("e2e_%v", value)
	userName := fmt.Sprintf("%v", value)
	email := fmt.Sprintf("%v@afl.com", value)
	return &Umpire{
		cli:      cli,
		Fullname: fullName,
		Username: userName,
		Email:    email,
	}
}
