package firestore

import (
	"context"
	"log"

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
