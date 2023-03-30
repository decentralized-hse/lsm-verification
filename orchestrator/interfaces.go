package orchestrator

type Orchestrator interface {
	SignNew() error
	// Returning last lseq, last hash
	ValidateFromLseq(lseqStart *string, hashLast* string) (string, string, error)
}
