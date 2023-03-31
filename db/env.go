package db

import (
	"crypto/rsa"
	"os"

	"lsm-verification/signature"
)

const (
	publicKeyEnvVariable  = "PUBLIC_KEY"
	privateKeyEnvVariable = "PRIVATE_KEY"
)

func loadPublicKey() (*rsa.PublicKey, error) {
	keyString := os.Getenv(publicKeyEnvVariable)
	if len(keyString) == 0 {
		return nil, ErrEmptyKey
	}

	return signature.LoadPublicKey(keyString)
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	keyString := os.Getenv(privateKeyEnvVariable)
	if len(keyString) == 0 {
		return nil, ErrEmptyKey
	}

	return signature.LoadPrivateKey(keyString)
}
