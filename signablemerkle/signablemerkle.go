package signablemerkle

import (
	"crypto/sha256"
	"errors"

	"github.com/cbergoon/merkletree"
	"github.com/decentralized-hse/lsm-verification/database"
)

type SignableDBItem database.DBItems_DbItem

func (s *SignableDBItem) CalculateHash() ([]byte, error) {
	stringRepr := (*database.DBItems_DbItem)(s).String()

	h := sha256.New()
	if _, err := h.Write([]byte(stringRepr)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (s *SignableDBItem) Equals(other merkletree.Content) (bool, error) {
	otherItem, ok := other.(*SignableDBItem)
	if !ok {
		return false, errors.New("value is not of type TestContent")
	}
	return (*database.DBItems_DbItem)(s).String() == (*database.DBItems_DbItem)(otherItem).String(), nil
}

type SignableMerkle merkletree.MerkleTree

func NewSignableMerkle(databaseItems []*database.DBItems_DbItem) (*SignableMerkle, error) {
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
