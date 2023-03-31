package db

import (
	"crypto/rsa"

	"lsm-verification/signature"
)

func loadPublicKey(keyString string) (*rsa.PublicKey, error) {
	if len(keyString) == 0 {
		return nil, ErrEmptyKey
	}

	return signature.LoadPublicKey(keyString)
}

func loadPrivateKey(keyString string) (*rsa.PrivateKey, error) {
	if len(keyString) == 0 {
		return nil, ErrEmptyKey
	}

	return signature.LoadPrivateKey(keyString)
}
