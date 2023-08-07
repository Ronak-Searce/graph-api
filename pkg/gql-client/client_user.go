package gqlclient

import (
	"fmt"
	"time"
)

type User struct {
	cli      *GQLClient
	Fullname string
	Username string
	Email    string
	Password string
}

func NewUser(cli *GQLClient) *User {
	value := time.Now().UTC().UnixMicro()
	fullName := fmt.Sprintf("e2e_%v", value)
	userName := fmt.Sprintf("%v", value)
	email := fmt.Sprintf("%v@afl.com", value)
	password := fmt.Sprintf("%v", value)
	return &User{
		cli:      cli,
		Fullname: fullName,
		Username: userName,
		Email:    email,
		Password: password,
	}
}
