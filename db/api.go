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
	eventsRequest := &proto.EventsRequest{
		ReplicaId: d.replicaId,
		Lseq:      startLseq,
		Limit:     &d.batchSize,
	}

	dbItemsObj, err := d.client.GetReplicaEvents(context.Background(), eventsRequest)
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

func (d *DbApi) getLastValue(key string) (*proto.Value, error) {
	replicaKey := &proto.ReplicaKey{
		Key:       key,
		ReplicaId: &d.replicaId,
	}
	return d.client.GetValue(context.Background(), replicaKey)
}

func (d *DbApi) ReadBatchValidated(lseqs []string) ([]models.ValidateItem, error) {
	result := make([]models.ValidateItem, 0, len(lseqs))

	for _, lseq := range lseqs {
		val, err := d.getLastValue(createValidationKey(lseq))
		if err != nil {
			return result, err
		}
		if val == nil {
			return result, nil // the rest is not validated, we skip it
		}

		hash, signed, err := splitHashAndSignature(val.Value)
		if err != nil {
			return result, err
		}

		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			return result, err
		}

		signatureBytes, err := hex.DecodeString(signed)
		if err != nil {
			return result, err
		}

		if err := signature.VerifySignature(signatureBytes, hashBytes, d.publicKey); err != nil {
			return result, err
		}

		result = append(
			result,
			models.ValidateItem{
				LseqItemValid: val.Lseq,
				Hash:          hash,
			},
		)
	}

	return result, nil
}

func (d *DbApi) GetLastValidated() (*models.ValidateItem, error) {
	validationValue, err := d.getLastValue(lastValidated)
	if err != nil || validationValue == nil {
		return nil, err
	}

	lastValidatedValue, err := d.getLastValue(createValidationKey(validationValue.Value))
	if err != nil {
		return nil, err
	}
	if lastValidatedValue == nil {
		return nil, ErrLastValidatedIsMissing
	}

	hash, _, err := splitHashAndSignature(lastValidatedValue.Value)
	if err != nil {
		return nil, err
	}
	result := &models.ValidateItem{
		Lseq:          &validationValue.Lseq,
		LseqItemValid: lastValidatedValue.Lseq,
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
