package services

import "chatonetoone/internals/auth"

type UserAuthenticator struct {
	auth auth.AuthenticationStrategy
}

func NewUserAuthenticator(strategy auth.AuthenticationStrategy) *UserAuthenticator {
	return &UserAuthenticator{auth: strategy}
}

func (ua *UserAuthenticator) Authenticate() (*auth.Result, error) {
	return ua.auth.Authenticate()
}
