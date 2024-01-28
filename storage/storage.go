package storage

import (
	"context"
	"io"
	"log"

	"cloud.google.com/go/storage"
	geojson "github.com/paulmach/go.geojson"
)

// ProvideStorage provides a storage client
func ProvideStorage() *storage.Client {
	client, err := storage.NewClient(context.TODO())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

var Options = ProvideStorage

func GetBucket() string {
	return "geotory-coldline"
}

func StoreObject(ctx context.Context, client *storage.Client, bucket, object string, data []byte) error {
	wc := client.Bucket(bucket).Object(object + ".json").NewWriter(ctx)
	wc.ContentType = "application/json"
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err := wc.Write(data); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func FetchObject(ctx context.Context, client *storage.Client, bucket, object string) ([]byte, error) {
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return io.ReadAll(rc)
}

type GeoJSONRespContext struct {
	Centroid []float64 `json:"centroid"`
	Zoom     int       `json:"zoom"`
}

type GeoJSONResp struct {
	Context GeoJSONRespContext         `json:"context"`
	GeoJSON *geojson.FeatureCollection `json:"geojson"`
}
