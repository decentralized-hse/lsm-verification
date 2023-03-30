package orchestrator

import "fmt"

var (
	ErrNoNewEntities = fmt.Errorf("No new entities found")
	ErrValidationFailed = fmt.Errorf("Validation failed")
	ErrBatchLenMismatch = fmt.Errorf("Batches length mismatch")
	ErrBadInput = fmt.Errorf("Bad input")
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
