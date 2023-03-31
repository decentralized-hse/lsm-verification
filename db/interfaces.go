package db

import "lsm-verification/models"

type DbState interface {
	CloseConnection()
	ReadBatch(startLseq *string) ([]models.DbItem, error)
	ReadBatchValidated(lseqs []string) ([]models.ValidateItem, error)
	GetLastValidated() (*models.ValidateItem, error)
	PutBatch(items []models.ValidateItem) error
}
