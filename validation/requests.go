package validation

import (
	"context"

	"github.com/decentralized-hse/lsm-verification/proto"
)

var HASH_KEY = "validation+hash"
var SIGNED_HASH_KEY = "validation+signedhash"
var LAST_LSEQ_KEY = "validation+lastlseq"

func getLastByKey(ctx context.Context, client proto.LSeqDatabaseClient, replicaId int32, key string) (string, error) {
	eventsRequest := &proto.EventsRequest{ReplicaId: replicaId, Key: &key}
	events, err := client.GetReplicaEvents(ctx, eventsRequest)
	if err != nil {
		return "", err
	}

	if events.ReplicaId != replicaId {
		return "", ErrWrongReplicaReturned
	}

	dbItems := events.Items
	if len(dbItems) == 0 {
		return "", ErrNoSuchKey
	}

	// TODO: check if order is guaranteed
	lastItem := dbItems[len(dbItems)-1]
	if lastItem.Key != HASH_KEY {
		return "", ErrWrongKeyReturned
	}

	return lastItem.Value, nil
}

func GetAll(
	ctx context.Context,
	client proto.LSeqDatabaseClient,
	replicaId int32,
	startLseq *string,
) ([]*proto.DBItems_DbItem, error) {
	eventsRequest := &proto.EventsRequest{ReplicaId: replicaId, Lseq: startLseq}
	events, err := client.GetReplicaEvents(ctx, eventsRequest)
	if err != nil {
		return nil, err
	}

	if events.ReplicaId != replicaId {
		return nil, ErrWrongReplicaReturned
	}

	filtered := make([]*proto.DBItems_DbItem, 0, len(events.Items))
	for _, item := range events.Items {
		if item.Key == HASH_KEY || item.Key == SIGNED_HASH_KEY || item.Key == LAST_LSEQ_KEY {
			continue
		}

		filtered = append(filtered, item)
	}

	if len(filtered) == 0 {
		return nil, ErrEmptyReplica
	}

	// TODO: check if order is guaranteed
	return filtered, nil
}

func GetLastHash(ctx context.Context, client proto.LSeqDatabaseClient, replicaId int32) (string, error) {
	return getLastByKey(ctx, client, replicaId, HASH_KEY)
}

func GetLastSignedHash(ctx context.Context, client proto.LSeqDatabaseClient, replicaId int32) (string, error) {
	return getLastByKey(ctx, client, replicaId, SIGNED_HASH_KEY)
}

func GetLastLseq(ctx context.Context, client proto.LSeqDatabaseClient, replicaId int32) (string, error) {
	return getLastByKey(ctx, client, replicaId, LAST_LSEQ_KEY)
}
