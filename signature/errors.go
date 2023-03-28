package signature

import "errors"

var ErrNoPEMBlock = errors.New("no PEM block found")
var ErrWrongKeyType = errors.New("wrong RSA key type")
var ErrUnableToParseKey = errors.New("unable to parse the key")
var ErrNoPassphrase = errors.New("no passphrase provided for an encrypted RSA private key")
