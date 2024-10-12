package utils

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func LoadPrivateKey(path string) error {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return err
	}
	return nil
}

func LoadPublicKey(path string) error {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return err
	}
	return nil
}
