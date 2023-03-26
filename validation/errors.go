package validation

import "errors"

var ErrNoSuchKey = errors.New("no such key exists in the replica")
var ErrEmptyReplica = errors.New("replica is empty")
var ErrWrongReplicaReturned = errors.New("server returned a different replica id from the requested one")
var ErrWrongKeyReturned = errors.New("server returned a different key from the requested one")
