package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders sets a small, safe set of hardening headers on every
// response. This is a pure JSON API (plus the Swagger UI page) with no
// user-supplied HTML rendering anywhere, so a strict Content-Security-Policy
// isn't added here - Swagger UI relies on inline scripts/styles that a real
// CSP would need to be tuned around, and getting that wrong risks breaking
// the docs page for no real security gain on JSON endpoints.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
