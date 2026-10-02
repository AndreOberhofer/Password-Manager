package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

var StaticKey []byte = []byte("32-byte-key-for-AES-256-!!!!!!!!")

func Encrypt(plainText string, key []byte) (string, error) {
	hashedKey := Hash32Bytes(key)
	block, err := aes.NewCipher(hashedKey[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12) // 12 bytes for GCM
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encodedCipher := base64.StdEncoding.EncodeToString(aesgcm.Seal(nonce, nonce, []byte(plainText), nil))

	return encodedCipher, nil
}

func Decrypt(ciphertext string, key []byte) (string, error) {
	hashedKey := Hash32Bytes(key)
	block, err := aes.NewCipher(hashedKey[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(decodedCiphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext is too short")
	}

	nonce, decodedCiphertext := decodedCiphertext[:nonceSize], decodedCiphertext[nonceSize:]

	plainText, err := aesgcm.Open(nil, nonce, decodedCiphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func Hash(plainText []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(plainText, 10)
}

func Hash32Bytes(plainText []byte) [32]byte {
	return sha256.Sum256(plainText)
}
