package auth

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestCognitoAuthorizerValidatesAccessTokenAndScopes(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	authorizer, err := NewCognitoAuthorizer(CognitoConfig{
		Region: "us-west-2", UserPoolID: "us-west-2_example", ClientID: "client-id",
	})
	if err != nil {
		t.Fatalf("create authorizer: %v", err)
	}
	authorizer.keys = map[string]*rsa.PublicKey{"key-1": &privateKey.PublicKey}
	authorizer.now = func() time.Time { return time.Unix(1_000, 0) }

	token := signTestJWT(t, privateKey, "key-1", map[string]any{
		"sub": "subject-1", "username": "operator",
		"client_id": "client-id", "iss": authorizer.issuer,
		"token_use": "access", "scope": "openid prism/proposals.write", "exp": 2_000,
	})
	identity, err := authorizer.Authorize(context.Background(), token, "prism/proposals.write")
	if err != nil {
		t.Fatalf("authorize token: %v", err)
	}
	if identity.Username != "operator" || identity.Subject != "subject-1" {
		t.Fatalf("unexpected identity: %+v", identity)
	}

	if _, err := authorizer.Authorize(context.Background(), token, "prism/admin.read"); !errors.Is(err, ErrInsufficientScope) {
		t.Fatalf("missing scope error = %v", err)
	}
}

func TestCognitoAuthorizerRejectsIDAndExpiredTokens(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	authorizer, _ := NewCognitoAuthorizer(CognitoConfig{
		Region: "us-west-2", UserPoolID: "us-west-2_example", ClientID: "client-id",
	})
	authorizer.keys = map[string]*rsa.PublicKey{"key-1": &privateKey.PublicKey}
	authorizer.now = func() time.Time { return time.Unix(1_000, 0) }

	for _, claims := range []map[string]any{
		{"sub": "subject-1", "client_id": "client-id", "iss": authorizer.issuer, "token_use": "id", "exp": 2_000},
		{"sub": "subject-1", "client_id": "client-id", "iss": authorizer.issuer, "token_use": "access", "exp": 999},
	} {
		token := signTestJWT(t, privateKey, "key-1", claims)
		if _, err := authorizer.Authorize(context.Background(), token); !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("invalid token error = %v", err)
		}
	}
}

func signTestJWT(t *testing.T, key *rsa.PrivateKey, keyID string, claims map[string]any) string {
	t.Helper()
	header, _ := json.Marshal(map[string]string{"alg": "RS256", "kid": keyID})
	payload, _ := json.Marshal(claims)
	unsigned := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	digest := sha256.Sum256([]byte(unsigned))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatalf("sign JWT: %v", err)
	}
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(signature)
}
