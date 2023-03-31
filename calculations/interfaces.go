package calculations

import (
	"errors"
	"lsm-verification/models"
)

type HashCalculator interface {
	CalculateBatch(items []models.DbItem, hashStart *string) ([]models.ValidateItem, error)
}

type hashCalculator struct{}

func (h *hashCalculator) CalculateBatch(items []models.DbItem, hashStart *string) ([]models.ValidateItem, error) {
	return nil, errors.New("Not implemented")
}

func CreateHashCalculator() HashCalculator {
	return &hashCalculator{}
}
