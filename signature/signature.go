package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
)

func Sign(dataHex string, privateKey *rsa.PrivateKey) (string, error) {
	dataBytes, err := hex.DecodeString(dataHex)
	if err != nil {
		return "", err
	}

	signatureBytes, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, dataBytes, nil)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signatureBytes), nil
}

func VerifySignature(signatureHex, dataHex string, publicKey *rsa.PublicKey) error {
	dataBytes, err := hex.DecodeString(dataHex)
	if err != nil {
		return err
	}

	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return err
	}

	return rsa.VerifyPSS(publicKey, crypto.SHA256, dataBytes, signatureBytes, nil)
}
