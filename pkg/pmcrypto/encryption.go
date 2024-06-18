package pmcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// 16 byte iv
// TODO: extract to config
var initializationVector = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Encrypt(target, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("new cipher: %v", err)
	}

	plainText := []byte(target)
	cfb := cipher.NewCFBEncrypter(block, initializationVector)

	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return encode(cipherText), nil
}

func decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func Decrypt(target, salt string) (string, error) {
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return "", fmt.Errorf("new cipher: %v", err)
	}

	cipherText, err := decode(target)
	if err != nil {
		return "", fmt.Errorf("decode: %v", err)
	}

	cfb := cipher.NewCFBDecrypter(block, initializationVector)

	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
