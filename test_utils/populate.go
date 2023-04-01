package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	"lsm-verification/models"
	"lsm-verification/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func put(client proto.LSeqDatabaseClient, key, value string) error {
	putRequest := &proto.PutRequest{
		Key:   key,
		Value: value,
	}

	log.Println("Requesting to append to the database")
	_, err := client.Put(context.Background(), putRequest)
	if err != nil {
		return err
	}
	log.Println("Appended to the database")

	return nil
}

func readBatch(d proto.LSeqDatabaseClient) ([]models.DbItem, error) {
	eventsRequest := &proto.EventsRequest{
		ReplicaId: 2,
	}

	log.Println("Requesting a DBItem batch from the database")
	dbItemsObj, err := d.GetReplicaEvents(context.Background(), eventsRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Received a DBItem batch from the database")

	dbItems := dbItemsObj.Items
	result := make([]models.DbItem, 0, len(dbItems))
	log.Println("Preprocessing the batch")
	for _, item := range dbItems {
		if item == nil {
			log.Fatalln("asldkjhasdljas")
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

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	addr, exists := os.LookupEnv("dbServerAddress")
	if !exists {
		log.Fatalln("Env variable 'dbServerAddress' not found")
	}

	log.Println("Dialing GRPC")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to start grpc connection", err)
	}

	log.Println("Creating a database client")
	client := proto.NewLSeqDatabaseClient(conn)

	var mode string
	fmt.Printf("Mode: ")
	fmt.Scanf("%s\n", &mode)

	if mode == "console" {
		for {
			var key string
			fmt.Printf("Key: ")
			fmt.Scanf("%s\n", &key)

			var value string
			fmt.Print("Value: ")
			fmt.Scanf("%s\n", &value)
			fmt.Println(key, value)
			put(client, key, value)
		}
	} else if mode == "read" {
		b, err := readBatch(client)
		if err != nil {
			log.Fatalln(err)
		}
		for i, v := range b {
			log.Println("Value number", i, ":", v)
		}
	} else {
		var count int
		fmt.Printf("count: ")
		fmt.Scanf("%d\n", &count)

		for i := 0; i < count; i += 1 {
			key := randSeq(5)
			value := randSeq(10)
			put(client, key, value)
			log.Println("Putted", key, value)
		}
		return
	}
}
