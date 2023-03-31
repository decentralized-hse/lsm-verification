package db

import "errors"

var ErrEmptyItem = errors.New("empty DB item pointer found")
var ErrLastValidatedIsMissing = errors.New("last validated lseq is not present in the database")
var ErrIncorrectValidationValue = errors.New("incorrect validation value, should be 'hash;signature'")
var ErrEmptyKey = errors.New("RSA key is empty")
var ErrNoKeys = errors.New("no RSA keys provided")
var ErrAddrNotSpecified = errors.New("server address is not specified")
var ErrReplicaIDNotSpecified = errors.New("replica ID is not specified")
var ErrPublicKeyEnvVarNotSpecified = errors.New("public key environment variable is not specified")
var ErrPrivateKeyEnvVarNotSpecified = errors.New("private key environment variable is not specified")
var ErrInvalidBatchSize = errors.New("batch size should be positive")
