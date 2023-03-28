package signature

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// Written with help from https://stackoverflow.com/questions/44230634/how-to-read-an-rsa-key-from-file

func loadPEM(keyPath string) (*pem.Block, error) {
	fileContents, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	decoded, _ := pem.Decode(fileContents)
	if decoded == nil {
		return nil, ErrNoPEMBlock
	}

	return decoded, nil
}

func LoadPrivateKey(rsaKeyPath string) (*rsa.PrivateKey, error) {
	pemBlock, err := loadPEM(rsaKeyPath)
	if err != nil {
		return nil, err
	}

	var parsedKey interface{}
	parsedKey, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
		if err != nil {
			return nil, err
		}
	}

	privKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrUnableToParseKey
	}

	return privKey, nil
}

func LoadPublicKey(rsaKeyPath string) (*rsa.PublicKey, error) {
	pemBlock, err := loadPEM(rsaKeyPath)
	if err != nil {
		return nil, err
	}

	parsedKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, ErrUnableToParseKey
	}

	pubKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrUnableToParseKey
	}

	return pubKey, nil
}
