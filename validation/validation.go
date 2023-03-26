package validation

import (
	"context"
	"crypto/rsa"

	"github.com/decentralized-hse/lsm-verification/proto"
	"github.com/decentralized-hse/lsm-verification/signablemerkle"
)

func Validate(ctx context.Context, client proto.LSeqDatabaseClient, replicaId int32, publicKey *rsa.PublicKey) error {
	lastHash, err := GetLastHash(ctx, client, replicaId)
	if err != nil {
		return err
	}

	lastHashSignature, err := GetLastSignedHash(ctx, client, replicaId)
	if err != nil {
		return err
	}

	dbItems, err := GetAllButValidation(ctx, client, replicaId, nil)
	if err != nil {
		return err
	}

	tree, err := signablemerkle.NewSignableMerkle(dbItems)
	if err != nil {
		return err
	}

	if err := tree.VerifyHash(lastHash); err != nil {
		return err
	}

	return tree.VerifyHashSignature(lastHashSignature, publicKey)
}
