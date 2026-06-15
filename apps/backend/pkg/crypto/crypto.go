package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

func DecodeEncryptionKey(secretKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(strings.Trim(strings.TrimSpace(secretKey), `"'`))
	if err != nil {
		return nil, errors.New("encryption key must be base64-encoded")
	}
	if len(key) != 32 {
		return nil, errors.New("encryption key must be exactly 32 bytes after base64 decoding")
	}
	return key, nil
}

func EncryptText(plainText, secretKey string) (string, error) {
	key, err := DecodeEncryptionKey(secretKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encryptedBytes := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func DecryptText(encryptedText, secretKey string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	key, err := DecodeEncryptionKey(secretKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("invalid encrypted value")
	}

	nonce := data[:gcm.NonceSize()]
	cipherText := data[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func MaskSecret(secret string) string {
	if len(secret) <= 10 {
		return "********"
	}

	return secret[:6] + "..." + secret[len(secret)-4:]
}
