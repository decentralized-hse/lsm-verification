package signablemerkle

import "errors"

var ErrWrongContentType = errors.New("value is not of type SignableDBItem")
var ErrInvalidHash = errors.New("invalid Merkle tree hash")
