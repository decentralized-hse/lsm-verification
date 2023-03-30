package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"lsm-verification/proto"
	"lsm-verification/signature"
	"lsm-verification/validation"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	VALIDATE      = "validate"
	HASH_AND_SIGN = "hash-and-sign"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatalf(
			"usage: %s {%s|%s} server_address replica_id {rsa_public_key|rsa_private_key}",
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
		key, err := signature.LoadPublicKey(rsaKeyPath)
		if err != nil {
			log.Fatal(err)
		}

		if err = validation.Validate(context.Background(), client, replicaId, key); err != nil {
			log.Fatal(err)
		}

		log.Println("validation was successful!")
	case HASH_AND_SIGN:
		key, err := signature.LoadPrivateKey(rsaKeyPath)
		if err != nil {
			log.Fatal(err)
		}

		if err = validation.HashAndSign(context.Background(), client, replicaId, key); err != nil {
			log.Fatal(err)
		}

		log.Println("hashing and signing were successful!")
	default:
		log.Fatalf("unknown command: %s", command)
	}
}
