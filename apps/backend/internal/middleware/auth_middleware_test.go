package middleware

import (
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func mustCIDR(t *testing.T, cidr string) *net.IPNet {
	t.Helper()
	_, network, err := net.ParseCIDR(cidr)
	require.NoError(t, err)
	return network
}

func TestRateLimiterStoreGetEvictsLeastRecentlySeenAtCapacity(t *testing.T) {
	t.Parallel()

	store := &RateLimiterStore{
		limiters: make(map[string]*clientLimiter, maxRateLimiterEntries),
		rate:     rate.Limit(1),
		burst:    1,
	}
	baseTime := time.Now().Add(-time.Hour)
	for i := 0; i < maxRateLimiterEntries; i++ {
		store.limiters["client-"+strconv.Itoa(i)] = &clientLimiter{
			limiter:  rate.NewLimiter(store.rate, store.burst),
			lastSeen: baseTime.Add(time.Duration(i) * time.Nanosecond),
		}
	}

	oldestLimiter := store.limiters["client-0"].limiter
	store.get("client-0")
	refreshedLastSeen := store.limiters["client-0"].lastSeen
	newLimiter := store.get("new-client")

	assert.Len(t, store.limiters, maxRateLimiterEntries)
	assert.Same(t, oldestLimiter, store.limiters["client-0"].limiter)
	assert.True(t, refreshedLastSeen.After(baseTime))
	assert.NotContains(t, store.limiters, "client-1")
	assert.Same(t, newLimiter, store.limiters["new-client"].limiter)
}

func TestClientIPForRateLimit_IgnoresXFFWithoutTrustedProxy(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		RemoteAddr: "203.0.113.10:54321",
		Header:     http.Header{"X-Forwarded-For": []string{"198.51.100.1"}},
	}

	ip := clientIPForRateLimit(req, nil)
	assert.Equal(t, "203.0.113.10", ip)
}

func TestClientIPForRateLimit_UsesXFFWhenRemoteIsTrusted(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		RemoteAddr: "10.0.0.5:443",
		Header:     http.Header{"X-Forwarded-For": []string{"198.51.100.1, 10.0.0.5"}},
	}
	trusted := []*net.IPNet{mustCIDR(t, "10.0.0.0/8")}

	ip := clientIPForRateLimit(req, trusted)
	assert.Equal(t, "198.51.100.1", ip)
}

func TestClientIPForRateLimit_RejectsLeftmostSpoofAppendedByTrustedProxy(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		RemoteAddr: "10.0.0.5:443",
		Header:     http.Header{"X-Forwarded-For": []string{"192.0.2.123, 198.51.100.20"}},
	}
	trusted := []*net.IPNet{mustCIDR(t, "10.0.0.0/8")}

	ip := clientIPForRateLimit(req, trusted)
	assert.Equal(t, "198.51.100.20", ip)
}

func TestClientIPForRateLimit_IgnoresXFFWhenRemoteUntrusted(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		RemoteAddr: "203.0.113.10:9999",
		Header:     http.Header{"X-Forwarded-For": []string{"198.51.100.1"}},
	}
	trusted := []*net.IPNet{mustCIDR(t, "10.0.0.0/8")}

	ip := clientIPForRateLimit(req, trusted)
	assert.Equal(t, "203.0.113.10", ip)
}

func TestClientIPForRateLimit_InvalidRemoteFallsBackUnknown(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		RemoteAddr: "not-an-ip",
		Header:     http.Header{"X-Forwarded-For": []string{"also-bad"}},
	}

	ip := clientIPForRateLimit(req, nil)
	assert.Equal(t, "unknown", ip)
}

func TestClientIPForRateLimit_TrustedButInvalidXFFUsesPeer(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		RemoteAddr: "10.0.0.5:443",
		Header:     http.Header{"X-Forwarded-For": []string{"not-an-ip"}},
	}
	trusted := []*net.IPNet{mustCIDR(t, "10.0.0.0/8")}

	ip := clientIPForRateLimit(req, trusted)
	assert.Equal(t, "10.0.0.5", ip)
}

func TestClientIPForRateLimit_TrustedOnlyXFFUsesPeer(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		RemoteAddr: "10.0.0.5:443",
		Header:     http.Header{"X-Forwarded-For": []string{"10.0.0.3, 10.0.0.4"}},
	}
	trusted := []*net.IPNet{mustCIDR(t, "10.0.0.0/8")}

	ip := clientIPForRateLimit(req, trusted)
	assert.Equal(t, "10.0.0.5", ip)
}

func TestParseIPFromRemoteAddr(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "127.0.0.1", parseIPFromRemoteAddr("127.0.0.1:8080").String())
	assert.Equal(t, "127.0.0.1", parseIPFromRemoteAddr("127.0.0.1").String())
	assert.Nil(t, parseIPFromRemoteAddr(""))
	assert.Nil(t, parseIPFromRemoteAddr("garbage"))
}

func TestFirstForwardedIP(t *testing.T) {
	t.Parallel()

	trusted := []*net.IPNet{mustCIDR(t, "10.0.0.0/8")}

	assert.Equal(t, "198.51.100.1", firstForwardedIP("198.51.100.1, 10.0.0.1", trusted).String())
	assert.Equal(t, "198.51.100.1", firstForwardedIP(" 198.51.100.1 ", trusted).String())
	assert.Nil(t, firstForwardedIP("10.0.0.1", trusted))
	assert.Nil(t, firstForwardedIP("", trusted))
	assert.Nil(t, firstForwardedIP("not-an-ip, still-bad", trusted))
}
