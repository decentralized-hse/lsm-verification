package calculations

import (
	"crypto/md5"
	"lsm-verification/db"
	"lsm-verification/models"
)

func CreateHashCalculator() HashCalculator {
	return &hashCalculator{}
}

func (h *hashCalculator) CalculateBatch(items []models.DbItem, hashStart *string) ([]models.ValidateItem, error) {
	if len(items) == 0 {
		return []models.ValidateItem{}, nil
	}

	result := make([]models.ValidateItem, 0, len(items))
	currentHash := ""
	if hashStart != nil {
		currentHash = *hashStart
	}

	for _, item := range items {
		currentHash = hashPrefixWithDbItem(&item, &currentHash)
		validatesItem := models.ValidateItem{
			Lseq:          &item.Lseq,
			LseqItemValid: db.CreateValidationKey(item.Lseq),
			Hash:          currentHash,
		}
		result = append(result, validatesItem)
	}

	return result, nil
}

func hashPrefixWithDbItem(item *models.DbItem, prefixHash *string) string {
	prefixWithlkv := *prefixHash + item.Lseq + item.Key + item.Value
	hash := md5.Sum([]byte(prefixWithlkv))
	return string(hash[:])
}
