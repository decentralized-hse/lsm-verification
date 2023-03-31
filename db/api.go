package db

import (
	"context"
	"crypto/rsa"
	"encoding/hex"
	"log"
	"lsm-verification/models"
	"lsm-verification/proto"
	"lsm-verification/signature"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultBatchSize = 100

type DbApi struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	replicaId  int32
	conn       *grpc.ClientConn
	client     proto.LSeqDatabaseClient
	batchSize  uint32
}

func NewDbApi(
	addr string,
	replicaId int32,
	privateKey *rsa.PrivateKey,
	publicKey *rsa.PublicKey,
	batchSize *uint32,
) DbApi {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := proto.NewLSeqDatabaseClient(conn)

	var finalBatchSize uint32 = defaultBatchSize
	if batchSize != nil {
		finalBatchSize = *batchSize
	}

	return DbApi{
		privateKey: privateKey,
		publicKey:  publicKey,
		replicaId:  replicaId,
		conn:       conn,
		client:     client,
		batchSize:  finalBatchSize,
	}
}

func (d *DbApi) CloseConnection() {
	d.conn.Close()
}

func (d *DbApi) ReadBatch(startLseq *string) ([]models.DbItem, error) {
	eventsRequest := proto.EventsRequest{
		ReplicaId: d.replicaId,
		Lseq:      startLseq,
		Limit:     &d.batchSize,
	}

	dbItemsObj, err := d.client.GetReplicaEvents(context.Background(), &eventsRequest)
	if err != nil {
		return nil, err
	}

	dbItems := dbItemsObj.Items
	result := make([]models.DbItem, 0, len(dbItems))
	for _, item := range dbItems {
		if item == nil {
			return nil, ErrEmptyItem
		}

		if isValidationKey(item.Key) {
			continue // we skip validation-specific keys
		}

		result = append(
			result,
			models.DbItem{
				Lseq:  item.Lseq,
				Key:   item.Key,
				Value: item.Value,
			},
		)
	}

	return result, nil
}

func (d *DbApi) ReadBatchValidated(lseqs []string) ([]models.ValidateItem, error) {

}

func (d *DbApi) getLastReplicaEvent(key string, noItemsErr error) (*proto.DBItems_DbItem, error) {
	eventsRequest := proto.EventsRequest{
		ReplicaId: d.replicaId,
		Key:       &key,
	}

	dbItemsObj, err := d.client.GetReplicaEvents(context.Background(), &eventsRequest)
	if err != nil {
		return nil, err
	}

	dbItems := dbItemsObj.Items
	if len(dbItems) == 0 {
		return nil, noItemsErr
	}

	lastItem := dbItems[len(dbItems)-1]
	if lastItem == nil {
		return nil, ErrEmptyItem
	}

	return lastItem, nil
}

func (d *DbApi) GetLastValidated() (*models.ValidateItem, error) {
	validationItem, err := d.getLastReplicaEvent(lastValidated, nil)
	if err != nil {
		return nil, err
	}
	if validationItem == nil {
		return nil, nil
	}

	lastValidatedItem, err := d.getLastReplicaEvent(
		createValidationKey(validationItem.Value),
		ErrLastValidatedIsMissing,
	)
	if err != nil {
		return nil, err
	}

	hash, _, err := splitHashAndSignature(lastValidatedItem.Value)
	if err != nil {
		return nil, err
	}
	result := &models.ValidateItem{
		Lseq:          &validationItem.Lseq,
		LseqItemValid: lastValidatedItem.Lseq,
		Hash:          hash,
	}
	return result, nil
}

func (d *DbApi) put(key, value string) error {
	putRequest := &proto.PutRequest{
		Key:   key,
		Value: value,
	}

	_, err := d.client.Put(context.Background(), putRequest)
	return err
}

func (d *DbApi) signAndPut(item models.ValidateItem) error {
	decoded, err := hex.DecodeString(item.Hash)
	if err != nil {
		return err
	}

	signedBytes, err := signature.Sign(decoded, d.privateKey)
	if err != nil {
		return err
	}

	signed := hex.EncodeToString(signedBytes)

	return d.put(item.LseqItemValid, joinHashAndSignature(item.Hash, signed))
}

func (d *DbApi) PutBatch(items []models.ValidateItem) error {
	for _, item := range items {
		if err := d.signAndPut(item); err != nil {
			return err
		}
	}

	return d.put(lastValidated, items[len(items)-1].LseqItemValid)
}
