package db

import (
	"crypto/rsa"
	"os"

	"lsm-verification/signature"
)

func loadPublicKey(envVariable string) (*rsa.PublicKey, error) {
	keyString := os.Getenv(envVariable)
	if len(keyString) == 0 {
		return nil, ErrEmptyKey
	}

	return signature.LoadPublicKey(keyString)
}

func loadPrivateKey(envVariable string) (*rsa.PrivateKey, error) {
	keyString := os.Getenv(envVariable)
	if len(keyString) == 0 {
		return nil, ErrEmptyKey
	}

	return signature.LoadPrivateKey(keyString)
}
