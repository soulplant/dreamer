package main

import (
	"golang.org/x/net/context"
	"github.com/soulplant/dreamer/api"
	"fmt"
	"errors"
)

// loginService handles login requests.
type loginService struct {}

func (l *loginService) Login(c context.Context, r *api.LoginRequest) (*api.LoginReply, error) {
	fmt.Println("loginService")
	if r.User == "" {
		return nil, errors.New("field 'user' required")
	}
	if r.Password == "" {
		return nil, errors.New("field 'password' required")
	}
	return &api.LoginReply{
		Token: "fake-token",
	}, nil
}
