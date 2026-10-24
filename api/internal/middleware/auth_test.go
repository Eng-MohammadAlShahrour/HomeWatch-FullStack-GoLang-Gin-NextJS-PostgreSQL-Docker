package middleware

import (
"net/http"
"net/http/httptest"
"testing"

"homelab/internal/auth"

"github.com/gin-gonic/gin"
)

func TestAuthRequiredMissingToken(t *testing.T) {
gin.SetMode(gin.TestMode)
w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)
c.Request = httptest.NewRequest("GET", "/", nil)

AuthRequired("secret")(c)

if w.Code != http.StatusUnauthorized {
t.Fatalf("expected 401, got %d", w.Code)
}
}

func TestAuthRequiredInvalidToken(t *testing.T) {
gin.SetMode(gin.TestMode)
w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)
c.Request = httptest.NewRequest("GET", "/", nil)
c.Request.Header.Set("Authorization", "Bearer invalid-token")

AuthRequired("secret")(c)

if w.Code != http.StatusUnauthorized {
t.Fatalf("expected 401, got %d", w.Code)
}
}

func TestAuthRequiredValidToken(t *testing.T) {
gin.SetMode(gin.TestMode)
w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)
c.Request = httptest.NewRequest("GET", "/", nil)

token, _ := auth.GenerateToken("secret", 7, "user@test.com")
c.Request.Header.Set("Authorization", "Bearer "+token)

AuthRequired("secret")(c)

if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d", w.Code)
}
if uid := c.GetUint("userID"); uid != 7 {
t.Fatalf("wrong userID: got %d", uid)
}
}