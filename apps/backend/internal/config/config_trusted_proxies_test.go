package config

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTrustedProxies_Empty(t *testing.T) {
	t.Parallel()
	nets, err := parseTrustedProxies("")
	require.NoError(t, err)
	assert.Nil(t, nets)
}

func TestParseTrustedProxies_IPsAndCIDRs(t *testing.T) {
	t.Parallel()
	nets, err := parseTrustedProxies("10.0.0.0/8, 127.0.0.1, 2001:db8::1")
	require.NoError(t, err)
	require.Len(t, nets, 3)

	assert.True(t, nets[0].Contains(net.ParseIP("10.1.2.3")))
	assert.True(t, nets[1].Contains(net.ParseIP("127.0.0.1")))
	assert.False(t, nets[1].Contains(net.ParseIP("127.0.0.2")))
	assert.True(t, nets[2].Contains(net.ParseIP("2001:db8::1")))
}

func TestParseTrustedProxies_Invalid(t *testing.T) {
	t.Parallel()
	_, err := parseTrustedProxies("not-a-cidr")
	require.Error(t, err)
}
