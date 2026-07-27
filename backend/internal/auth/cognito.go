package auth

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

type CognitoConfig struct {
	Region     string
	UserPoolID string
	ClientID   string
	HTTPClient *http.Client
}

type CognitoAuthorizer struct {
	issuer   string
	clientID string
	jwksURL  string
	client   *http.Client
	now      func() time.Time
	mu       sync.RWMutex
	keys     map[string]*rsa.PublicKey
}

type cognitoClaims struct {
	Subject   string `json:"sub"`
	Username  string `json:"username"`
	ClientID  string `json:"client_id"`
	Issuer    string `json:"iss"`
	TokenUse  string `json:"token_use"`
	Scope     string `json:"scope"`
	ExpiresAt int64  `json:"exp"`
}

func NewCognitoAuthorizer(cfg CognitoConfig) (*CognitoAuthorizer, error) {
	if strings.TrimSpace(cfg.Region) == "" || strings.TrimSpace(cfg.UserPoolID) == "" || strings.TrimSpace(cfg.ClientID) == "" {
		return nil, fmt.Errorf("Cognito region, user pool ID, and client ID are required")
	}
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", cfg.Region, cfg.UserPoolID)
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &CognitoAuthorizer{issuer: issuer, clientID: cfg.ClientID, jwksURL: issuer + "/.well-known/jwks.json", client: client, now: time.Now}, nil
}

func (a *CognitoAuthorizer) Authorize(ctx context.Context, token string, requiredScopes ...string) (Identity, error) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 {
		return Identity{}, ErrInvalidToken
	}
	var header struct {
		Algorithm string `json:"alg"`
		KeyID     string `json:"kid"`
	}
	if err := decodeJWTPart(parts[0], &header); err != nil || header.Algorithm != "RS256" || header.KeyID == "" {
		return Identity{}, ErrInvalidToken
	}
	var claims cognitoClaims
	if err := decodeJWTPart(parts[1], &claims); err != nil {
		return Identity{}, ErrInvalidToken
	}
	key, err := a.key(ctx, header.KeyID)
	if err != nil {
		return Identity{}, ErrInvalidToken
	}
	digest := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || rsa.VerifyPKCS1v15(key, crypto.SHA256, digest[:], signature) != nil {
		return Identity{}, ErrInvalidToken
	}
	if claims.Issuer != a.issuer || claims.ClientID != a.clientID || claims.TokenUse != "access" || claims.ExpiresAt <= a.now().Unix() || claims.Subject == "" {
		return Identity{}, ErrInvalidToken
	}
	scopes := strings.Fields(claims.Scope)
	for _, required := range requiredScopes {
		if !contains(scopes, required) {
			return Identity{}, ErrInsufficientScope
		}
	}
	username := claims.Username
	if username == "" {
		username = claims.Subject
	}
	return Identity{Username: username, Subject: claims.Subject, Scopes: scopes}, nil
}

func (a *CognitoAuthorizer) key(ctx context.Context, keyID string) (*rsa.PublicKey, error) {
	a.mu.RLock()
	key := a.keys[keyID]
	a.mu.RUnlock()
	if key != nil {
		return key, nil
	}
	if err := a.refreshKeys(ctx); err != nil {
		return nil, err
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	key = a.keys[keyID]
	if key == nil {
		return nil, fmt.Errorf("Cognito signing key %q not found", keyID)
	}
	return key, nil
}

func (a *CognitoAuthorizer) refreshKeys(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.jwksURL, nil)
	if err != nil {
		return err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Cognito JWKS returned %s", resp.Status)
	}
	var document struct {
		Keys []struct {
			KeyID    string `json:"kid"`
			KeyType  string `json:"kty"`
			Use      string `json:"use"`
			Modulus  string `json:"n"`
			Exponent string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return err
	}
	keys := make(map[string]*rsa.PublicKey)
	for _, item := range document.Keys {
		if item.KeyType != "RSA" || item.Use != "sig" {
			continue
		}
		n, errN := base64.RawURLEncoding.DecodeString(item.Modulus)
		e, errE := base64.RawURLEncoding.DecodeString(item.Exponent)
		if errN != nil || errE != nil || len(e) == 0 {
			continue
		}
		exponent := 0
		for _, b := range e {
			exponent = exponent<<8 + int(b)
		}
		key := &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: exponent}
		if _, err := x509.MarshalPKIXPublicKey(key); err == nil {
			keys[item.KeyID] = key
		}
	}
	if len(keys) == 0 {
		return fmt.Errorf("Cognito JWKS contains no usable signing keys")
	}
	a.mu.Lock()
	a.keys = keys
	a.mu.Unlock()
	return nil
}

func decodeJWTPart(value string, target any) error {
	data, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
