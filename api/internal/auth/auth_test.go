package auth

import (
"testing"
"time"
)

func TestHashPassword(t *testing.T) {
hash, err := HashPassword("secret123")
if err != nil {
t.Fatalf("hash failed: %v", err)
}
if hash == "secret123" {
t.Fatal("password was not hashed")
}
if len(hash) < 50 {
t.Fatalf("hash too short: %d", len(hash))
}
}

func TestCheckPasswordCorrect(t *testing.T) {
hash, _ := HashPassword("secret123")
if !CheckPassword(hash, "secret123") {
t.Fatal("correct password rejected")
}
}

func TestCheckPasswordWrong(t *testing.T) {
hash, _ := HashPassword("secret123")
if CheckPassword(hash, "wrong") {
t.Fatal("wrong password accepted")
}
}

func TestGenerateAndParseToken(t *testing.T) {
secret := "test-secret"
token, err := GenerateToken(secret, 42, "test@example.com")
if err != nil {
t.Fatalf("token generation failed: %v", err)
}
if token == "" {
t.Fatal("empty token")
}

claims, err := ParseToken(secret, token)
if err != nil {
t.Fatalf("parse failed: %v", err)
}
if claims["email"] != "test@example.com" {
t.Fatalf("wrong email: %v", claims["email"])
}
if uint(claims["sub"].(float64)) != 42 {
t.Fatalf("wrong sub: %v", claims["sub"])
}
}

func TestParseTokenInvalid(t *testing.T) {
_, err := ParseToken("secret", "invalid-token")
if err == nil {
t.Fatal("expected error for invalid token")
}
}

func TestParseTokenWrongSecret(t *testing.T) {
token, _ := GenerateToken("secret1", 1, "a@b.com")
_, err := ParseToken("secret2", token)
if err == nil {
t.Fatal("expected error for wrong secret")
}
}

func TestTokenExpiry(t *testing.T) {
token, _ := GenerateToken("secret", 1, "a@b.com")
claims, _ := ParseToken("secret", token)
exp := int64(claims["exp"].(float64))
now := time.Now().Unix()
if exp <= now {
t.Fatal("token already expired")
}
if exp > now+25*3600 {
t.Fatal("token expiry too far")
}
}