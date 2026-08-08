package middleware

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustCIDR(t *testing.T, cidr string) *net.IPNet {
	t.Helper()
	_, network, err := net.ParseCIDR(cidr)
	require.NoError(t, err)
	return network
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

func TestParseIPFromRemoteAddr(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "127.0.0.1", parseIPFromRemoteAddr("127.0.0.1:8080").String())
	assert.Equal(t, "127.0.0.1", parseIPFromRemoteAddr("127.0.0.1").String())
	assert.Nil(t, parseIPFromRemoteAddr(""))
	assert.Nil(t, parseIPFromRemoteAddr("garbage"))
}

func TestFirstForwardedIP(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "198.51.100.1", firstForwardedIP("198.51.100.1, 10.0.0.1").String())
	assert.Equal(t, "198.51.100.1", firstForwardedIP(" 198.51.100.1 ").String())
	assert.Nil(t, firstForwardedIP(""))
	assert.Nil(t, firstForwardedIP("not-an-ip, still-bad"))
}
