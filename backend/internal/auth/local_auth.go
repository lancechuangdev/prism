package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type LocalConfig struct {
	AdminUsername string
	AdminPassword string
	TokenSecret   string
	TokenTTL      time.Duration
}

type LocalAuthenticator struct {
	adminUsername string
	adminPassword string
	tokenSecret   []byte
	tokenTTL      time.Duration
	now           func() time.Time
	sessions      SessionStore
}

type tokenPayload struct {
	Username  string `json:"username"`
	ExpiresAt int64  `json:"expiresAt"`
	Nonce     string `json:"nonce"`
}

func NewLocalAuthenticator(cfg LocalConfig, sessions SessionStore) *LocalAuthenticator {
	return &LocalAuthenticator{
		adminUsername: cfg.AdminUsername,
		adminPassword: cfg.AdminPassword,
		tokenSecret:   []byte(cfg.TokenSecret),
		tokenTTL:      cfg.TokenTTL,
		now:           time.Now,
		sessions:      sessions,
	}
}

func (s *LocalAuthenticator) Login(ctx context.Context, username string, password string) (string, error) {
	if username != s.adminUsername || password != s.adminPassword {
		return "", ErrInvalidCredentials
	}

	expiresAt := s.now().UTC().Add(s.tokenTTL)
	payload := tokenPayload{
		Username:  username,
		ExpiresAt: expiresAt.Unix(),
		Nonce:     randomNonce(),
	}

	token, err := s.sign(payload)
	if err != nil {
		return "", err
	}

	if err := s.sessions.Set(ctx, sessionKey(token), username, s.tokenTTL); err != nil {
		return "", fmt.Errorf("store session: %w", err)
	}
	return token, nil
}

func (s *LocalAuthenticator) Logout(ctx context.Context, rawToken string) error {
	token := strings.TrimSpace(rawToken)
	if token == "" {
		return nil
	}
	return s.sessions.Delete(ctx, sessionKey(token))
}

func (s *LocalAuthenticator) Authenticate(ctx context.Context, rawToken string) (string, error) {
	token := strings.TrimSpace(rawToken)
	if token == "" {
		return "", ErrInvalidToken
	}

	payload, err := s.verify(token)
	if err != nil {
		return "", err
	}

	now := s.now().UTC()
	if now.Unix() >= payload.ExpiresAt {
		return "", ErrInvalidToken
	}

	username, err := s.sessions.Get(ctx, sessionKey(token))
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return "", ErrInvalidToken
		}
		return "", fmt.Errorf("load session: %w", err)
	}
	if username != payload.Username {
		return "", ErrInvalidToken
	}
	return username, nil
}

func sessionKey(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("auth:session:%x", hash[:])
}

func randomNonce() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func (s *LocalAuthenticator) sign(payload tokenPayload) (string, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signature, err := signValue(encodedPayload, s.tokenSecret)
	if err != nil {
		return "", err
	}
	return encodedPayload + "." + signature, nil
}

func (s *LocalAuthenticator) verify(token string) (tokenPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return tokenPayload{}, ErrInvalidToken
	}

	payloadStr := parts[0]
	expectedSignature, err := signValue(payloadStr, s.tokenSecret)
	if err != nil {
		return tokenPayload{}, ErrInvalidToken
	}
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[1])) {
		return tokenPayload{}, ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadStr)
	if err != nil {
		return tokenPayload{}, ErrInvalidToken
	}

	payload := tokenPayload{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return tokenPayload{}, ErrInvalidToken
	}

	return payload, nil
}

func signValue(value string, secret []byte) (string, error) {
	mac := hmac.New(sha256.New, secret)
	if _, err := mac.Write([]byte(value)); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}
