package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	validKey   = "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoxMjM0NTY="
	shortKey   = "YWJjZGVmZ2hpamtsbW5vcA=="
	invalidB64 = "not-a-valid-base64!!!"
)

// ── DecodeEncryptionKey ────────────────────────────────────────────────

func TestDecodeEncryptionKey_Success(t *testing.T) {
	key, err := DecodeEncryptionKey(validKey)

	require.NoError(t, err)
	assert.Len(t, key, 32)
}

func TestDecodeEncryptionKey_StripsQuotes(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"double quotes", `"` + validKey + `"`},
		{"single quotes", `'` + validKey + `'`},
		{"leading/trailing spaces", `  ` + validKey + `  `},
		{"mixed", `  "` + validKey + `"  `},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := DecodeEncryptionKey(tt.input)
			require.NoError(t, err)
			assert.Len(t, key, 32)
		})
	}
}

func TestDecodeEncryptionKey_InvalidBase64(t *testing.T) {
	_, err := DecodeEncryptionKey(invalidB64)

	assert.ErrorContains(t, err, "encryption key must be base64-encoded")
}

func TestDecodeEncryptionKey_WrongLength(t *testing.T) {
	_, err := DecodeEncryptionKey(shortKey)

	assert.ErrorContains(t, err, "encryption key must be exactly 32 bytes")
}

// ── EncryptText / DecryptText ─────────────────────────────────────────

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := "hello world"

	ciphertext, err := EncryptText(plaintext, validKey)
	require.NoError(t, err)
	require.NotEmpty(t, ciphertext)
	require.NotEqual(t, plaintext, ciphertext)

	decrypted, err := DecryptText(ciphertext, validKey)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptDecrypt_SpecialCharacters(t *testing.T) {
	plaintext := "Hello, 世界! 🎉\n\t\"'$`\\"

	ciphertext, err := EncryptText(plaintext, validKey)
	require.NoError(t, err)

	decrypted, err := DecryptText(ciphertext, validKey)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptDecrypt_UniqueCiphertext(t *testing.T) {
	c1, err := EncryptText("same data", validKey)
	require.NoError(t, err)

	c2, err := EncryptText("same data", validKey)
	require.NoError(t, err)

	assert.NotEqual(t, c1, c2, "each encryption should produce a unique nonce")
}

func TestEncryptDecrypt_WrongKey(t *testing.T) {
	ciphertext, err := EncryptText("secret message", validKey)
	require.NoError(t, err)

	otherKey := "bWFuaG9iYXZlbGxhbmd1YWdlZHdvcmxkZG9n"
	_, err = DecryptText(ciphertext, otherKey)

	assert.Error(t, err)
}

func TestEncryptDecrypt_InvalidCiphertext(t *testing.T) {
	_, err := DecryptText("not-valid-base64!!!", validKey)
	assert.Error(t, err)
}

func TestEncryptDecrypt_InvalidKey(t *testing.T) {
	_, err := EncryptText("hello", invalidB64)
	assert.ErrorContains(t, err, "encryption key must be base64-encoded")

	_, err = DecryptText("dGVzdA==", invalidB64)
	assert.ErrorContains(t, err, "encryption key must be base64-encoded")
}

func TestEncryptDecrypt_EmptyString(t *testing.T) {
	ciphertext, err := EncryptText("", validKey)
	require.NoError(t, err)

	decrypted, err := DecryptText(ciphertext, validKey)
	require.NoError(t, err)
	assert.Equal(t, "", decrypted)
}

// ── MaskSecret ────────────────────────────────────────────────────────

func TestMaskSecret_Long(t *testing.T) {
	result := MaskSecret("abcdefghijklmnopqrstu")
	assert.Equal(t, "abcdef...rstu", result)
}

func TestMaskSecret_Short(t *testing.T) {
	assert.Equal(t, "********", MaskSecret("abcdefghij"))
	assert.Equal(t, "********", MaskSecret("abc"))
	assert.Equal(t, "********", MaskSecret(""))
}

func TestMaskSecret_Edge(t *testing.T) {
	assert.Equal(t, "abcdef...stuv", MaskSecret("abcdefghijstuv"))
	assert.Equal(t, "abcdef...hijk", MaskSecret("abcdefghijk"))
}
