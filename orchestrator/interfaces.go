package orchestrator

import "errors"

var (
	ErrNoNewEntities = errors.New("No new entities found")
	ErrValidationFailed = errors.New("Validation failed")
	ErrBatchLenMismatch = errors.New("Batches length mismatch")
	ErrBadInput = errors.New("Bad input")
)

type Orchestrator interface {
	SignNew() error

	/* 
	 * Returning
	 * last lseq validated or lseq which is failed validation
	 * last hash if we have something to validate
	 */
	ValidateFromLseq(lseqStart *string, hashLast* string) (*string, *string, error)
}
