package service

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := "ghp_test_token_1234567890abcdef"

	encrypted, err := Encrypt(plaintext, key)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, plaintext, encrypted)

	decrypted, err := Decrypt(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncrypt_WrongKeyLength(t *testing.T) {
	key := make([]byte, 16) // Too short
	_, err := Encrypt("test", key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	_, _ = rand.Read(key1)
	_, _ = rand.Read(key2)

	encrypted, err := Encrypt("secret data", key1)
	require.NoError(t, err)

	_, err = Decrypt(encrypted, key2)
	assert.Error(t, err)
}

func TestDecrypt_InvalidCiphertext(t *testing.T) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)

	// Not base64
	_, err := Decrypt("not-valid-base64!!!", key)
	assert.Error(t, err)

	// Valid base64 but invalid ciphertext
	_, err = Decrypt("dGVzdA==", key)
	assert.Error(t, err)
}

func TestDecrypt_WrongKeyLength(t *testing.T) {
	key := make([]byte, 16) // Too short
	_, err := Decrypt("dGVzdA==", key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")
}

func TestEncryptDecrypt_EmptyString(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encrypted, err := Encrypt("", key)
	require.NoError(t, err)

	decrypted, err := Decrypt(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, "", decrypted)
}
