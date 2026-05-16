package auth

type CredentialValidator interface {
	Validate(input string, source string) error
}
