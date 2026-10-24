package handlers

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"github.com/gin-gonic/gin"
)

func TestRegisterInvalidEmail(t *testing.T) {
gin.SetMode(gin.TestMode)
w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)

body, _ := json.Marshal(map[string]string{"email": "not-an-email", "password": "secret123"})
c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
c.Request.Header.Set("Content-Type", "application/json")

Register(c)

if w.Code != http.StatusBadRequest {
t.Fatalf("expected 400, got %d", w.Code)
}
}

func TestRegisterShortPassword(t *testing.T) {
gin.SetMode(gin.TestMode)
w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)

body, _ := json.Marshal(map[string]string{"email": "a@b.com", "password": "123"})
c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
c.Request.Header.Set("Content-Type", "application/json")

Register(c)

if w.Code != http.StatusBadRequest {
t.Fatalf("expected 400, got %d", w.Code)
}
}