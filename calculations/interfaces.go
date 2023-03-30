package calculations

import "lsm-verification/models"

type HashCalculator interface {
	CalculateBatch(items []models.DbItem, hashStart *string) ([]models.ValidateItem, error)
}
