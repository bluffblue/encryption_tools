package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

type PasswordEncryptionService struct {
	iterations int
	keySize    int
}

func NewPasswordEncryptionService() *PasswordEncryptionService {
	return &PasswordEncryptionService{
		iterations: 100000,
		keySize:    32,
	}
}

func (s *PasswordEncryptionService) DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, s.iterations, s.keySize, sha256.New)
}

func (s *PasswordEncryptionService) Encrypt(plaintext, password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := s.DeriveKey(password, salt)

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

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	ciphertext = append(salt, ciphertext...)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *PasswordEncryptionService) Decrypt(encryptedText, password string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	salt := ciphertext[:16]
	ciphertext = ciphertext[16:]

	key := s.DeriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", err
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
