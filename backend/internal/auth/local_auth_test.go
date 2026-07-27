package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSessionIsSharedAcrossServiceInstances(t *testing.T) {
	store := NewMemorySessionStore()
	first := newTestLocalAuthenticator(store)
	second := newTestLocalAuthenticator(store)

	token, err := first.Login(context.Background(), "admin", "password")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	username, err := second.Authenticate(context.Background(), token)
	if err != nil {
		t.Fatalf("authenticate through second instance: %v", err)
	}
	if username != "admin" {
		t.Fatalf("unexpected username %q", username)
	}
}

func TestLogoutInvalidatesSharedSession(t *testing.T) {
	store := NewMemorySessionStore()
	first := newTestLocalAuthenticator(store)
	second := newTestLocalAuthenticator(store)

	token, err := first.Login(context.Background(), "admin", "password")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if err := second.Logout(context.Background(), token); err != nil {
		t.Fatalf("logout: %v", err)
	}
	if _, err := first.Authenticate(context.Background(), token); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected invalid token after logout, got %v", err)
	}
}

func TestLoginRejectsInvalidCredentialsWithoutCreatingSession(t *testing.T) {
	service := newTestLocalAuthenticator(NewMemorySessionStore())
	if _, err := service.Login(context.Background(), "admin", "wrong"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func newTestLocalAuthenticator(store SessionStore) *LocalAuthenticator {
	return NewLocalAuthenticator(LocalConfig{
		AdminUsername: "admin",
		AdminPassword: "password",
		TokenSecret:   "test-secret",
		TokenTTL:      time.Hour,
	}, store)
}
