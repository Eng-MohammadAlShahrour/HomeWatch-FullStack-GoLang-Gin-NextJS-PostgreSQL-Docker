package middleware

import (
"strconv"
"time"

"homelab/internal/metrics"

"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
if c.Request.URL.Path == "/metrics" {
c.Next()
return
}
start := time.Now()
c.Next()
duration := time.Since(start).Seconds()
path := c.FullPath()
if path == "" {
path = "unmatched"
}
metrics.HTTPRequestsTotal.WithLabelValues(c.Request.Method, path, strconv.Itoa(c.Writer.Status())).Inc()
metrics.HTTPDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
}
}