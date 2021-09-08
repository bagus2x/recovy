package middleware

import "github.com/bagus2x/recovy/auth"

type Middleware struct {
	authService auth.Service
}

func NewMiddleware(authService auth.Service) *Middleware {
	return &Middleware{
		authService: authService,
	}
}
