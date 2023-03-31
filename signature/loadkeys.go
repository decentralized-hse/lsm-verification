package signature

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// Written with help from https://stackoverflow.com/questions/44230634/how-to-read-an-rsa-key-from-file

func loadPEM(rsaKeyContents string) (*pem.Block, error) {
	decoded, _ := pem.Decode([]byte(rsaKeyContents))
	if decoded == nil {
		return nil, ErrNoPEMBlock
	}

	return decoded, nil
}

func LoadPrivateKey(rsaKeyContents string) (*rsa.PrivateKey, error) {
	pemBlock, err := loadPEM(rsaKeyContents)
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

func LoadPublicKey(rsaKeyContents string) (*rsa.PublicKey, error) {
	pemBlock, err := loadPEM(rsaKeyContents)
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
