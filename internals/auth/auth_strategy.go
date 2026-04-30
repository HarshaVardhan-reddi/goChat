package auth

type Result struct {
	IsSuccess  bool
	AuthResult map[string]string
	UserData   map[string]string
}

type AuthenticationStrategy interface {
	Authenticate() (*Result, error)
}
