package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// visitor tracks one client IP's request budget.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// ipRateLimiter is a simple in-memory, per-process, per-IP limiter.
//
// This is intentionally lightweight (no Redis/external store), which is
// the right tradeoff for this app's current single-instance deployment
// model. If NileStay is ever run as multiple replicas behind a load
// balancer, each instance enforces its own independent limit rather than
// a shared one - a real distributed limiter (e.g. Redis-backed) would be
// needed at that point. That infrastructure doesn't exist in this project
// today, so this is flagged here rather than pretended away.
type ipRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     rate.Limit
	burst    int
}

func newIPRateLimiter(r rate.Limit, burst int) *ipRateLimiter {
	l := &ipRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    burst,
	}
	go l.cleanupLoop()
	return l
}

func (l *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	v, exists := l.visitors[ip]
	if !exists {
		lim := rate.NewLimiter(l.rate, l.burst)
		l.visitors[ip] = &visitor{limiter: lim, lastSeen: time.Now()}
		return lim
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// cleanupLoop prevents unbounded memory growth from tracking every
// distinct IP forever - entries idle for more than 5 minutes are dropped.
func (l *ipRateLimiter) cleanupLoop() {
	for {
		time.Sleep(time.Minute)
		l.mu.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

// AuthRateLimiter limits login/register attempts per client IP - 5 requests
// per minute with a burst of 5, generous enough for normal retries after a
// typo but tight enough to blunt credential-stuffing/brute-force attempts
// against /auth/login and enumeration/spam against /auth/register.
func AuthRateLimiter() gin.HandlerFunc {
	limiter := newIPRateLimiter(rate.Every(12*time.Second), 5)

	return func(c *gin.Context) {
		if !limiter.getLimiter(c.ClientIP()).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "too many requests, please try again shortly",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
