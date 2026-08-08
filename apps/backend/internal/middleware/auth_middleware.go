package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/pkg/jwt"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterStore struct {
	mu             sync.Mutex
	limiters       map[string]*clientLimiter
	rate           rate.Limit // request per second
	burst          int
	logger         *zap.Logger
	trustedProxies []*net.IPNet
}

func NewRateLimiterStore(r rate.Limit, burst int, logger *zap.Logger, trustedProxies []*net.IPNet) *RateLimiterStore {
	store := &RateLimiterStore{
		limiters:       make(map[string]*clientLimiter),
		rate:           r,
		burst:          burst,
		logger:         logger,
		trustedProxies: trustedProxies,
	}
	// Background cleanup: remove entries not seen in 10 minutes
	go store.cleanupLoop()
	return store
}

func (s *RateLimiterStore) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)

	for range ticker.C {
		s.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for key, entry := range s.limiters {
			if entry.lastSeen.Before(cutoff) {
				delete(s.limiters, key)
			}
		}
		s.mu.Unlock()
	}
}

func (s *RateLimiterStore) get(key string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.limiters[key]

	if !exists {
		entry = &clientLimiter{
			limiter: rate.NewLimiter(s.rate, s.burst),
		}
		s.limiters[key] = entry
	}

	entry.lastSeen = time.Now()

	return entry.limiter
}

func IPRateLimiter(store *RateLimiterStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := clientIPForRateLimit(c.Request, store.trustedProxies)
		if ip == "unknown" && store.logger != nil {
			store.logger.Warn(
				"rate limit using unknown client IP",
				zap.String("remote_addr", c.Request.RemoteAddr),
				zap.String("x_forwarded_for", c.GetHeader("X-Forwarded-For")),
			)
		}

		if !store.get(ip).Allow() {
			if store.logger != nil {
				store.logger.Warn(
					"rate limit exceeded",
					zap.String("ip", ip),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
				)
			}
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests from your IP",
			})
			return
		}
		c.Next()
	}
}

// clientIPForRateLimit returns a validated client IP for rate limiting.
// X-Forwarded-For is only trusted when RemoteAddr belongs to a configured trusted proxy.
func clientIPForRateLimit(r *http.Request, trustedProxies []*net.IPNet) string {
	peer := parseIPFromRemoteAddr(r.RemoteAddr)
	if peer != nil && isTrustedProxy(peer, trustedProxies) {
		if forwarded := firstForwardedIP(r.Header.Get("X-Forwarded-For"), trustedProxies); forwarded != nil {
			return forwarded.String()
		}
	}
	if peer != nil {
		return peer.String()
	}
	return "unknown"
}

func parseIPFromRemoteAddr(remoteAddr string) net.IP {
	if remoteAddr == "" {
		return nil
	}
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// Bare IP without port (or already stripped).
		host = remoteAddr
	}
	// Strip zone identifier from IPv6 (e.g. fe80::1%lo0).
	if i := strings.IndexByte(host, '%'); i >= 0 {
		host = host[:i]
	}
	return net.ParseIP(host)
}

func firstForwardedIP(xff string, trustedProxies []*net.IPNet) net.IP {
	if xff == "" {
		return nil
	}
	parts := strings.Split(xff, ",")
	for i := len(parts) - 1; i >= 0; i-- {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}
		if ip := parseIPFromRemoteAddr(part); ip != nil && !isTrustedProxy(ip, trustedProxies) {
			return ip
		}
	}
	return nil
}

func isTrustedProxy(ip net.IP, trusted []*net.IPNet) bool {
	if ip == nil || len(trusted) == 0 {
		return false
	}
	for _, network := range trusted {
		if network != nil && network.Contains(ip) {
			return true
		}
	}
	return false
}
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := jwt.ValidateJwt(tokenString, cfg)
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
