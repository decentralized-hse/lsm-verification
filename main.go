package main

import (
	"context"
	"crypto/rsa"
	"log"
	"os"
	"strconv"

	"github.com/decentralized-hse/lsm-verification/proto"
	"github.com/decentralized-hse/lsm-verification/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	VALIDATE      = "validate"
	HASH_AND_SIGN = "hash-and-sign"
)

// TODO: finish key loading
func loadPrivateKey(rsaKeyPath string) (*rsa.PrivateKey, error) {
	panic("not implemented")
}

func loadPublicKey(rsaKeyPath string) (*rsa.PublicKey, error) {
	panic("not implemented")
}

func main() {
	if len(os.Args) != 5 {
		log.Fatalf(
			"Usage: %s {%s|%s} server_address replica_id {rsa_public_key|rsa_private_key}",
			os.Args[0],
			VALIDATE,
			HASH_AND_SIGN,
		)
	}

	command := os.Args[1]
	addr := os.Args[2]
	replicaIdInt, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	replicaId := int32(replicaIdInt)
	rsaKeyPath := os.Args[4]

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewLSeqDatabaseClient(conn)

	switch command {
	case VALIDATE:
		key, err := loadPublicKey(rsaKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		err = validation.Validate(context.Background(), client, replicaId, key)
	case HASH_AND_SIGN:
		key, err := loadPrivateKey(rsaKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		err = validation.HashAndSign(context.Background(), client, replicaId, key)
	default:
		log.Fatalf("unknown command: %s", command)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Operation successful!")
}
