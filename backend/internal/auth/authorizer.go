package auth

import (
	"context"
	"errors"
)

var ErrInsufficientScope = errors.New("insufficient scope")

type Identity struct {
	Username string
	Subject  string
	Scopes   []string
}

type Authorizer interface {
	Authorize(ctx context.Context, token string, requiredScopes ...string) (Identity, error)
}

type LocalAuthorizer struct {
	service *LocalAuthenticator
}

func NewLocalAuthorizer(service *LocalAuthenticator) *LocalAuthorizer {
	return &LocalAuthorizer{service: service}
}

func (a *LocalAuthorizer) Authorize(ctx context.Context, token string, _ ...string) (Identity, error) {
	username, err := a.service.Authenticate(ctx, token)
	if err != nil {
		return Identity{}, err
	}
	return Identity{Username: username, Subject: username}, nil
}
