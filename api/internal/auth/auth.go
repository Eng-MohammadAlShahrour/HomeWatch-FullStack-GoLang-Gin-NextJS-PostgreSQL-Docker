package auth

import (
"time"

"github.com/golang-jwt/jwt/v5"
"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) (string, error) {
b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
return string(b), err
}

func CheckPassword(hash, p string) bool {
return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}

func GenerateToken(secret string, userID uint, email string) (string, error) {
claims := jwt.MapClaims{
"sub":   userID,
"email": email,
"exp":   time.Now().Add(24 * time.Hour).Unix(),
"iat":   time.Now().Unix(),
}
t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return t.SignedString([]byte(secret))
}

func ParseToken(secret, tokenStr string) (jwt.MapClaims, error) {
token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
return []byte(secret), nil
})
if err != nil {
return nil, err
}
claims, ok := token.Claims.(jwt.MapClaims)
if !ok || !token.Valid {
return nil, jwt.ErrTokenInvalidClaims
}
return claims, nil
}