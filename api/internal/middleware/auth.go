package middleware

import (
"net/http"
"strings"

"homelab/internal/auth"

"github.com/gin-gonic/gin"
)

func AuthRequired(secret string) gin.HandlerFunc {
return func(c *gin.Context) {
h := c.GetHeader("Authorization")
if !strings.HasPrefix(h, "Bearer ") {
c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
return
}
tokenStr := strings.TrimPrefix(h, "Bearer ")
claims, err := auth.ParseToken(secret, tokenStr)
if err != nil {
c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
return
}
c.Set("userID", uint(claims["sub"].(float64)))
c.Set("email", claims["email"])
c.Next()
}
}