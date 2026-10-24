package middleware

import (
"net/http"
"sync"
"time"

"github.com/gin-gonic/gin"
)

type bucket struct {
tokens float64
last   time.Time
}

type RateLimiter struct {
mu       sync.Mutex
buckets  map[string]*bucket
capacity float64
rate     float64
}

func NewRateLimiter(capacity int, refillPerSec float64) *RateLimiter {
rl := &RateLimiter{
buckets:  map[string]*bucket{},
capacity: float64(capacity),
rate:     refillPerSec,
}
go rl.cleanup()
return rl
}

func (rl *RateLimiter) cleanup() {
t := time.NewTicker(10 * time.Minute)
defer t.Stop()
for range t.C {
rl.mu.Lock()
cutoff := time.Now().Add(-30 * time.Minute)
for k, b := range rl.buckets {
if b.last.Before(cutoff) {
delete(rl.buckets, k)
}
}
rl.mu.Unlock()
}
}

func (rl *RateLimiter) Allow(key string) bool {
rl.mu.Lock()
defer rl.mu.Unlock()
now := time.Now()
b, ok := rl.buckets[key]
if !ok {
rl.buckets[key] = &bucket{tokens: rl.capacity - 1, last: now}
return true
}
elapsed := now.Sub(b.last).Seconds()
b.tokens += elapsed * rl.rate
if b.tokens > rl.capacity {
b.tokens = rl.capacity
}
b.last = now
if b.tokens < 1 {
return false
}
b.tokens -= 1
return true
}

func RateLimit(rl *RateLimiter) gin.HandlerFunc {
return func(c *gin.Context) {
key := c.ClientIP()
if !rl.Allow(key) {
c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
return
}
c.Next()
}
}