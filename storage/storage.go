package storage

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// ProvideStorage provides a storage client
func ProvideStorage() *storage.Client {
	client, err := storage.NewClient(context.TODO(), option.WithCredentialsFile("/Users/mager/Code/bluedot/credentials.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

var Options = ProvideStorage
