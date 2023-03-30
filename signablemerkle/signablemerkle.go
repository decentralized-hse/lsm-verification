package signablemerkle

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"

	"lsm-verification/proto"
	"lsm-verification/signature"

	"github.com/cbergoon/merkletree"
)

type SignableDBItem proto.DBItems_DbItem

func (s *SignableDBItem) CalculateHash() ([]byte, error) {
	stringRepr := (*proto.DBItems_DbItem)(s).String()

	h := sha256.New()
	if _, err := h.Write([]byte(stringRepr)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (s *SignableDBItem) Equals(other merkletree.Content) (bool, error) {
	otherItem, ok := other.(*SignableDBItem)
	if !ok {
		return false, ErrWrongContentType
	}
	return (*proto.DBItems_DbItem)(s).String() == (*proto.DBItems_DbItem)(otherItem).String(), nil
}

type SignableMerkle merkletree.MerkleTree

func NewSignableMerkle(databaseItems []*proto.DBItems_DbItem) (*SignableMerkle, error) {
	convertedDatabaseItems := make([]merkletree.Content, 0, len(databaseItems))
	for _, dbItem := range databaseItems {
		convertedDatabaseItems = append(convertedDatabaseItems, (*SignableDBItem)(dbItem))
	}

	tree, err := merkletree.NewTree(convertedDatabaseItems)
	if err != nil {
		return nil, err
	}

	return (*SignableMerkle)(tree), nil
}

func (s *SignableMerkle) GetHash() string {
	return hex.EncodeToString(s.Root.Hash)
}

func (s *SignableMerkle) GetSignedHash(privateKey *rsa.PrivateKey) (string, error) {
	signature, err := signature.Sign(s.Root.Hash, privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signature), nil
}

func (s *SignableMerkle) VerifyHash(hash string) error {
	if hash != s.GetHash() {
		return ErrInvalidHash
	}

	return nil
}

func (s *SignableMerkle) VerifyHashSignature(signatureStr string, publicKey *rsa.PublicKey) error {
	signatureBytes, err := hex.DecodeString(signatureStr)
	if err != nil {
		return err
	}
	return signature.VerifySignature(signatureBytes, s.Root.Hash, publicKey)
}
