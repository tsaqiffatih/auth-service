package utils

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

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

func GenerateJWT(encryptedID string, email string) (string, error) {

	claims := jwt.MapClaims{
		"id":    encryptedID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // token valid for 3 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}
