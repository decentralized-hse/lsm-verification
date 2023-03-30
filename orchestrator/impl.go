package orchestrator

import (
	"log"

	"lsm-verification/models"
	"lsm-verification/db"
	"lsm-verification/calculations"
)

type orchestrator struct {
	db db.DbState
	calculator calculations.HashCalculator
}

func (o *orchestrator) SignNew() error {
	lastValidated, err := o.db.GetLastValidated()
	if err != nil {
		return err
	}
	log.Println("Got last validated")

	var batch []models.DbItem
	if lastValidated != nil {
		batch, err = o.db.ReadBatch(&lastValidated.LseqItemValid)
	} else {
		batch, err = o.db.ReadBatch(nil)
	}
	if err != nil {
		return err
	}
	if len(batch) == 0 {
		log.Println("Batch is empty")
		return ErrNoNewEntities
	}
	log.Println("Got batch")

	var calculatedBatch []models.ValidateItem
	if lastValidated != nil {
		calculatedBatch, err = o.calculator.CalculateBatch(batch, &lastValidated.Hash)
	} else {
		calculatedBatch, err = o.calculator.CalculateBatch(batch, nil)
	}
	if err != nil {
		return err
	}
	if (len(batch) != len(calculatedBatch)) {
		return ErrBatchLenMismatch
	}

	log.Println("Putting validated batch")
	return o.db.PutBatch(calculatedBatch)
}


func (o *orchestrator) ValidateFromLseq(lseqStart *string, hashLast* string) (*string, *string, error) {
	if (lseqStart == nil && hashLast != nil) || (lseqStart != nil && hashLast == nil) {
		return nil, nil, ErrBadInput
	}

	batch, err := o.db.ReadBatch(lseqStart)
	if err != nil {
		return nil, nil, err
	}
	if len(batch) == 0 {
		log.Println("Batch is empty, validation done")
		return nil, nil, ErrNoNewEntities
	}
	log.Println("Got batch")

	calculatedBatch, err := o.calculator.CalculateBatch(batch, hashLast)
	if err != nil {
		return nil, nil, err
	}
	if (len(batch) != len(calculatedBatch)) {
		return nil, nil, ErrBatchLenMismatch
	}
	log.Println("Calculated batch")

	lseqs := []string{}
	for _, item := range batch {
		lseqs = append(lseqs, item.Lseq)
	}
	validBatch, err := o.db.ReadBatchValidated(lseqs)
	if err != nil {
		return nil, nil, err
	}
	log.Println("Got validated batch")

	if (len(calculatedBatch) > len(validBatch)) {
		log.Println("Batches arent equal, validated to last valid")
	} else if (len(calculatedBatch) < len(validBatch)) {
		return nil, nil, ErrBatchLenMismatch
	}

	for idx, validItem := range validBatch {
		if (validItem.Hash != calculatedBatch[idx].Hash) {
			log.Println("Batch not valid on lseq:", validItem.LseqItemValid)
			return &validItem.LseqItemValid, nil, ErrValidationFailed
		}
	}

	log.Println("Batch is valid")
	lastItem := &validBatch[len(validBatch)-1]
	if (len(calculatedBatch) > len(validBatch)) {
		return &lastItem.LseqItemValid, nil, nil
	}
	return &lastItem.LseqItemValid, &lastItem.Hash, nil
}

func CreateOrchestrator(db db.DbState, calculator calculations.HashCalculator) Orchestrator {
	return &orchestrator{
		db: db,
		calculator: calculator,
	}
}
