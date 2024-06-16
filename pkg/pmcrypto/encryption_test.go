package pmcrypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		target := "test input"
		secret := "1234567890123456"

		encryptedText, err := Encrypt(target, secret)
		require.NoError(t, err, "Encrypt should not return an error for valid inputs")
		require.NotEmpty(t, encryptedText, "Encrypt should return a non-empty string")
		require.NotEqual(t, target, encryptedText)
	})

	t.Run("error_invalid_secret_length", func(t *testing.T) {
		target := "test input"
		invalidSecret := "short"

		_, err := Encrypt(target, invalidSecret)
		require.Error(t, err, "Encrypt should return an error for an invalid secret length")
	})
}

func TestDecrypt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		target := "test input"
		secret := "1234567890123456"

		encryptedText, err := Encrypt(target, secret)
		require.NoError(t, err, "Encrypt should not return an error for valid inputs")

		decryptedText, err := Decrypt(encryptedText, secret)
		require.NoError(t, err, "Decrypt should not return an error for valid inputs")
		require.Equal(t, target, decryptedText, "Decrypted text should match the original text")
	})

	t.Run("error_invalid_secret_length", func(t *testing.T) {
		text := "test input"
		invalidSecret := "short"

		_, err := Decrypt(text, invalidSecret)
		require.Error(t, err, "Decrypt should return an error for an invalid secret length")
	})
}

func BenchmarkEncrypt(b *testing.B) {
	text := "hello world"
	secret := "1234567890123456"

	for i := 0; i < b.N; i++ {
		Encrypt(text, secret)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	text := "hello world"
	secret := "1234567890123456"

	encryptedText, err := Encrypt(text, secret)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		Decrypt(encryptedText, secret)
	}
}
