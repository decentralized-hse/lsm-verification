package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

func Sign(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, data, nil)
}

func VerifySignature(signature, data []byte, publicKey *rsa.PublicKey) error {
	return rsa.VerifyPSS(publicKey, crypto.SHA256, data, signature, nil)
}
