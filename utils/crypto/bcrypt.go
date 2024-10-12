package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(hashedPassword string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

func formatKey(key string) []byte {
	keyBytes := []byte(key)
	keyLength := len(keyBytes)

	if keyLength < 16 {
		// Pad key to 16 bytes
		padding := make([]byte, 16-keyLength)
		keyBytes = append(keyBytes, padding...)
	} else if keyLength > 16 && keyLength < 24 {
		// Pad key to 24 bytes
		padding := make([]byte, 24-keyLength)
		keyBytes = append(keyBytes, padding...)
	} else if keyLength > 24 && keyLength < 32 {
		// Pad key to 32 bytes
		padding := make([]byte, 32-keyLength)
		keyBytes = append(keyBytes, padding...)
	} else if keyLength > 32 {
		// Trim key to 32 bytes
		keyBytes = keyBytes[:32]
	}

	// Now the keyBytes is guaranteed to be either 16, 24, or 32 bytes
	return keyBytes
}

func EncryptString(plaintext string, key string) (string, error) {
	// Format the key to ensure it's of valid length
	keyBytes := formatKey(key)
	plaintextBytes := []byte(plaintext)

	// Print the length of the key to confirm it
	fmt.Println("Key length:", len(keyBytes))

	// Create a new AES cipher using the provided key
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("NewCipher failed: %v", err)
	}

	fmt.Println("2")
	// Generate a new IV (Initialization Vector)
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("IV generation failed: %v", err)
	}

	fmt.Println("3")
	// Encrypt the data using CFB mode
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)

	fmt.Println("4")

	// Encode the ciphertext as a base64 string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
