package firestore

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
)

// ProvideFirestore provides a firestore client
func ProvideFirestore() *firestore.Client {
	projectID := "geotory"

	client, err := firestore.NewClient(context.TODO(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

var Options = ProvideFirestore

type Dataset struct {
	Image string  `json:"image" firestore:"image"`
	Types []int64 `json:"types" firestore:"types"`
}

type Feature struct {
	Type int `json:"type" firestore:"type"`
}

type Property struct {
	Dataset string      `json:"dataset" firestore:"dataset"`
	Data    interface{} `json:"data" firestore:"data"`
}

type Geometry struct {
	Dataset string      `json:"dataset" firestore:"dataset"`
	Data    interface{} `json:"data" firestore:"data"`
}

const (
	// DatasetTypeValueUnknown is an unknown dataset type
	DatasetTypeValueUnknown int = iota
	// DatasetTypeGeopackage is a geopackage dataset type
	DatasetTypeGeopackage
	// DatasetTypeGeojson is a geojson dataset type
	DatasetTypeGeojson
	// DatasetTypeShapefile is a shapefile dataset type
	DatasetTypeShapefile

	// DatasetTypeNameGeopackage is the name of the geopackage dataset type
	DatasetTypeNameGeopackage = "geopackage"
	// DatasetTypeNameGeojson is the name of the geojson dataset type
	DatasetTypeNameGeojson = "geojson"
	// DatasetTypeNameShapefile is the name of the shapefile dataset type
	DatasetTypeNameShapefile = "shapefile"
)

func DatasetTypeNameToValue(typeName string) int {
	switch typeName {
	case DatasetTypeNameGeopackage:
		return DatasetTypeGeopackage
	case DatasetTypeNameGeojson:
		return DatasetTypeGeojson
	case DatasetTypeNameShapefile:
		return DatasetTypeShapefile
	default:
		return DatasetTypeValueUnknown
	}
}

func DatasetTypeValueToName(typeValue int) string {
	switch typeValue {
	case DatasetTypeGeopackage:
		return DatasetTypeNameGeopackage
	case DatasetTypeGeojson:
		return DatasetTypeNameGeojson
	case DatasetTypeShapefile:
		return DatasetTypeNameShapefile
	default:
		return ""
	}
}

func GenerateDocumentID(name string) string {
	// Generate a random number between 0 and 2^31
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	randomNumber := uint32(binary.BigEndian.Uint32(randomBytes))

	// Generate a timestamp in milliseconds
	timestamp := uint32(time.Now().UnixNano() / 1000000)

	// Combine the timestamp and random number to generate the document ID
	prefix := fmt.Sprintf("%010d%08x", timestamp, randomNumber)
	return prefix + name
}
