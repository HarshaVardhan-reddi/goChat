package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordValidator struct{}

func (pv *PasswordValidator) Validate(input string, source string) error {
	if input == "" {
		return errors.New("auth: input password is missing")
	}
	if source == "" {
		return errors.New("auth: source password hash is missing")
	}

	err := bcrypt.CompareHashAndPassword([]byte(source), []byte(input))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("auth: invalid password")
		}
		return err
	}
	return nil
}
