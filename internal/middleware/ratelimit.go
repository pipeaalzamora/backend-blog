package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
}

var limiter = &rateLimiter{requests: make(map[string][]time.Time)}

func RateLimit(maxPerSecond int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		limiter.mu.Lock()
		times := limiter.requests[ip]
		var recent []time.Time
		for _, t := range times {
			if now.Sub(t) < time.Second {
				recent = append(recent, t)
			}
		}
		if len(recent) >= maxPerSecond {
			limiter.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		limiter.requests[ip] = append(recent, now)
		limiter.mu.Unlock()
		c.Next()
	}
}
