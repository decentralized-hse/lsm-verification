package db

import "strings"

const validationPrefix = "_v_"

const lastValidated = "_v_asd"

func isValidationKey(key string) bool {
	return strings.HasPrefix(key, validationPrefix)
}

func createValidationKey(key string) string {
	return validationPrefix + key
}

func splitHashAndSignature(joined string) (string, string, error) {
	split := strings.Split(joined, ";")
	if len(split) != 2 {
		return "", "", ErrIncorrectValidationValue
	}

	return split[0], split[1], nil
}

func joinHashAndSignature(hash, signature string) string {
	return hash + ";" + signature
}
