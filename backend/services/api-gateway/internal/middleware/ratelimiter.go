package middleware

import (
	"net/http"
	"sync"
	"time"

	"api-gateway/pkg/response"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.RWMutex
)

func RateLimiter(requestsPerMinute int) gin.HandlerFunc {
	go func() {
		for {
			time.Sleep(3 * time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{
				lastSeen: time.Now(),
				count:    1,
			}
			mu.Unlock()
			c.Next()
			return
		}

		if time.Since(v.lastSeen) > time.Minute {
			v.count = 1
			v.lastSeen = time.Now()
			mu.Unlock()
			c.Next()
			return
		}

		if v.count >= requestsPerMinute {
			mu.Unlock()
			response.Error(c, http.StatusTooManyRequests, "Rate limit exceeded")
			c.Abort()
			return
		}

		v.count++
		v.lastSeen = time.Now()
		mu.Unlock()

		c.Next()
	}
}
