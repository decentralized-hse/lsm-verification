package db

import (
	"context"
	"crypto/rsa"
	"log"
	"lsm-verification/config"
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

func CreateDbState(cfg config.Config) (DbState, error) {
	addr, ok := cfg.Env.Db.ServerAddress.(string)
	if !ok {
		return nil, ErrAddrNotSpecified
	}

	replicaId, ok := cfg.Env.Db.ReplicaID.(int)
	if !ok {
		return nil, ErrReplicaIDNotSpecified
	}

	var batchSize *uint32
	if bs, ok := cfg.Db.BatchSize.(int); ok {
		if bs <= 0 {
			return nil, ErrInvalidBatchSize
		}

		*batchSize = uint32(bs)
	}

	publicKeyEnvVar, ok := cfg.Env.Rsa.PublicKey.(string)
	if !ok {
		return nil, ErrPublicKeyEnvVarNotSpecified
	}

	privateKeyEnvVar, ok := cfg.Env.Rsa.PrivateKey.(string)
	if !ok {
		return nil, ErrPrivateKeyEnvVarNotSpecified
	}

	dbApi, err := CreateDbApi(
		addr,
		int32(replicaId),
		batchSize,
		publicKeyEnvVar,
		privateKeyEnvVar,
	)
	if err != nil {
		return nil, err
	}

	return dbApi, nil
}

func CreateDbApi(
	addr string,
	replicaId int32,
	batchSize *uint32,
	publicKeyEnvVariable string,
	privateKeyEnvVariable string,
) (*DbApi, error) {
	log.Println("Dialing GRPC")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	log.Println("Creating a database connection")
	client := proto.NewLSeqDatabaseClient(conn)

	var finalBatchSize uint32 = defaultBatchSize
	if batchSize != nil && *batchSize != 0 {
		finalBatchSize = *batchSize
	}
	log.Printf("Set the database batch size as %d\n", finalBatchSize)

	log.Println("Trying to load the public key")
	publicKey, err := loadPublicKey(publicKeyEnvVariable)
	if err != nil {
		if err == ErrEmptyKey {
			log.Println("Warning: public key is not set, can only certify history")
		} else {
			return nil, err
		}
	}

	log.Println("Trying to load the private key")
	privateKey, err := loadPrivateKey(privateKeyEnvVariable)
	if err != nil {
		if err == ErrEmptyKey {
			log.Println("Warning: private key is not set, can only verify history")
		} else {
			return nil, err
		}
	}

	if publicKey == nil && privateKey == nil {
		return nil, ErrNoKeys
	}

	return &DbApi{
		privateKey: privateKey,
		publicKey:  publicKey,
		replicaId:  replicaId,
		conn:       conn,
		client:     client,
		batchSize:  finalBatchSize,
	}, nil
}

func (d *DbApi) CloseConnection() {
	log.Println("Closing the database connection")
	d.conn.Close()
}

func (d *DbApi) ReadBatch(startLseq *string) ([]models.DbItem, error) {
	eventsRequest := &proto.EventsRequest{
		ReplicaId: d.replicaId,
		Lseq:      startLseq,
		Limit:     &d.batchSize,
	}

	log.Println("Requesting a DBItem batch from the database")
	dbItemsObj, err := d.client.GetReplicaEvents(context.Background(), eventsRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Received a DBItem batch from the database")

	dbItems := dbItemsObj.Items
	result := make([]models.DbItem, 0, len(dbItems))
	log.Println("Preprocessing the batch")
	for _, item := range dbItems {
		if item == nil {
			return nil, ErrEmptyItem
		}

		if isValidationKey(item.Key) {
			log.Println("Skipping a validation-specific key")
			continue
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
	log.Println("Finished preprocessing the batch")

	return result, nil
}

func (d *DbApi) getLastValue(key string) (*proto.Value, error) {
	replicaKey := &proto.ReplicaKey{
		Key:       key,
		ReplicaId: &d.replicaId,
	}

	log.Println("Requesting the last value based on a key from the database")
	val, err := d.client.GetValue(context.Background(), replicaKey)
	if err != nil {
		return nil, err
	}
	log.Println("Received the last value based on a key from the database")

	return val, nil
}

func (d *DbApi) ReadBatchValidated(lseqs []string) ([]models.ValidateItem, error) {
	result := make([]models.ValidateItem, 0, len(lseqs))

	log.Println("Validating lseqs")
	for _, lseq := range lseqs {
		val, err := d.getLastValue(createValidationKey(lseq))
		if err != nil {
			return result, err
		}
		if val == nil {
			log.Println("Got an unvalidated lseq, skipping the rest")
			return result, nil
		}
		log.Println("Loaded a validated lseq")

		log.Println("Splitting its value into the hash and the signature")
		hash, signed, err := splitHashAndSignature(val.Value)
		if err != nil {
			return result, err
		}

		log.Println("Verifying the signature")
		if err := signature.VerifySignature(signed, hash, d.publicKey); err != nil {
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
	log.Println("Successfully validated all lseqs")

	return result, nil
}

func (d *DbApi) GetLastValidated() (*models.ValidateItem, error) {
	validationValue, err := d.getLastValue(lastValidated)
	if err != nil || validationValue == nil {
		return nil, err
	}
	log.Println("Loaded the last validated lseq object")

	lastValidatedValue, err := d.getLastValue(createValidationKey(validationValue.Value))
	if err != nil {
		return nil, err
	}
	if lastValidatedValue == nil {
		return nil, ErrLastValidatedIsMissing
	}
	log.Println("Loaded the hash and signature object for the last validated lseq")

	log.Println("Splitting its value into the hash and the signature")
	hash, signed, err := splitHashAndSignature(lastValidatedValue.Value)
	if err != nil {
		return nil, err
	}

	log.Println("Verifying the signature")
	if err := signature.VerifySignature(signed, hash, d.publicKey); err != nil {
		return nil, err
	}

	result := &models.ValidateItem{
		Lseq:          &validationValue.Lseq,
		LseqItemValid: lastValidatedValue.Lseq,
		Hash:          hash,
	}
	log.Println("Constructed the last validated item")

	return result, nil
}

func (d *DbApi) put(key, value string) error {
	putRequest := &proto.PutRequest{
		Key:   key,
		Value: value,
	}

	log.Println("Requesting to append to the database")
	_, err := d.client.Put(context.Background(), putRequest)
	if err != nil {
		return err
	}
	log.Println("Appended to the database")

	return nil
}

func (d *DbApi) signAndPut(item models.ValidateItem) error {
	log.Println("Signing the hash")
	signed, err := signature.Sign(item.Hash, d.privateKey)
	if err != nil {
		return err
	}

	return d.put(item.LseqItemValid, joinHashAndSignature(item.Hash, signed))
}

func (d *DbApi) PutBatch(items []models.ValidateItem) error {
	log.Println("Appending a batch to the database")
	for _, item := range items {
		if err := d.signAndPut(item); err != nil {
			return err
		}
	}

	log.Println("Updating the last validated lseq in the database")
	return d.put(lastValidated, items[len(items)-1].LseqItemValid)
}
