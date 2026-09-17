package jwt

import (
	"testing"
	"time"

	realjwt "github.com/golang-jwt/jwt/v5"
)

func setupTestEnv(t *testing.T) {
	t.Helper()
	t.Setenv("JWT_SECRET", "test-secret-do-not-use-in-production")
	t.Setenv("JWT_EXPIRE", "1h")
}

func TestGenerateAndParseToken_RoundTrip(t *testing.T) {
	setupTestEnv(t)

	token, err := GenerateToken("user-123", "admin")
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("expected no error parsing a freshly generated token, got %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("expected UserID 'user-123', got %q", claims.UserID)
	}
	if claims.Role != "admin" {
		t.Errorf("expected Role 'admin', got %q", claims.Role)
	}
}

func TestParseToken_RejectsExpiredToken(t *testing.T) {
	setupTestEnv(t)

	claims := Claims{
		UserID: "user-123",
		Role:   "customer",
		RegisteredClaims: realjwt.RegisteredClaims{
			ExpiresAt: realjwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  realjwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := realjwt.NewWithClaims(realjwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("test-secret-do-not-use-in-production"))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	if _, err := ParseToken(signed); err == nil {
		t.Fatal("expected an error parsing an expired token, got nil")
	}
}

func TestParseToken_RejectsWrongSigningAlgorithm(t *testing.T) {
	setupTestEnv(t)

	claims := Claims{
		UserID: "user-123",
		Role:   "admin",
		RegisteredClaims: realjwt.RegisteredClaims{
			ExpiresAt: realjwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	// Signed with HS384 instead of the HS256 this app issues and expects -
	// WithValidMethods([]string{"HS256"}) in ParseToken must reject this
	// even though the signature itself is otherwise valid for that key.
	token := realjwt.NewWithClaims(realjwt.SigningMethodHS384, claims)
	signed, err := token.SignedString([]byte("test-secret-do-not-use-in-production"))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	if _, err := ParseToken(signed); err == nil {
		t.Fatal("expected an error parsing a token signed with a non-HS256 algorithm, got nil")
	}
}

func TestParseToken_RejectsTamperedToken(t *testing.T) {
	setupTestEnv(t)

	token, err := GenerateToken("user-123", "customer")
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	tampered := token[:len(token)-2] + "xx"

	if _, err := ParseToken(tampered); err == nil {
		t.Fatal("expected an error parsing a tampered token, got nil")
	}
}

func TestParseToken_RejectsMalformedToken(t *testing.T) {
	setupTestEnv(t)

	if _, err := ParseToken("not-a-real-jwt"); err == nil {
		t.Fatal("expected an error parsing a malformed token, got nil")
	}
}

func TestParseToken_RejectsWrongSecret(t *testing.T) {
	setupTestEnv(t)

	token, err := GenerateToken("user-123", "customer")
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	t.Setenv("JWT_SECRET", "a-completely-different-secret")

	if _, err := ParseToken(token); err == nil {
		t.Fatal("expected an error parsing a token signed with a different secret, got nil")
	}
}
