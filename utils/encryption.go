package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
)

func EncryptID(id uint) (string, error) {
	plainText := []byte(strconv.FormatUint(uint64(id), 10))

	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plainText, nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptID(encryptedID string) (uint, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encryptedID)
	if err != nil {
		return 0, err
	}

	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherText, nil)
	if err != nil {
		return 0, err
	}

	// Convert plainText back to uint
	id, err := strconv.ParseUint(string(plainText), 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
