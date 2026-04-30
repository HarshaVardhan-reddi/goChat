package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrMissingInput    = errors.New("auth: input password is missing")
	ErrMissingSource   = errors.New("auth: source password hash is missing")
	ErrInvalidPassword = errors.New("auth: invalid password")
)

type AuthPassword struct {
	sourcePasswordHash []byte
	inputPassword      []byte
}

func NewAuthPassword(sourcePasswordHash, inputPassword string) *AuthPassword {
	return &AuthPassword{
		sourcePasswordHash: []byte(sourcePasswordHash),
		inputPassword:      []byte(inputPassword),
	}
}

func (atp *AuthPassword) Authenticate() (*Result, error) {
	if err := atp.validateArgs(); err != nil {
		return nil, err
	}

	err := bcrypt.CompareHashAndPassword(atp.sourcePasswordHash, atp.inputPassword)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidPassword
		}
		return nil, err
	}
	return &Result{IsSuccess: true}, nil
}

func (atp *AuthPassword) validateArgs() error {
	if len(atp.inputPassword) == 0 {
		return ErrMissingInput
	}
	if len(atp.sourcePasswordHash) == 0 {
		return ErrMissingSource
	}
	return nil
}
