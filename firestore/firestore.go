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
	Image string `json:"image"`
}

type DatasetTypeName string
type DatasetTypeValue int

const (
	// DatasetTypeValueUnknown is an unknown dataset type
	DatasetTypeValueUnknown DatasetTypeValue = iota
	// DatasetTypeGeopackage is a geopackage dataset type
	DatasetTypeGeopackage
	// DatasetTypeGeojson is a geojson dataset type
	DatasetTypeGeojson
	// DatasetTypeShapefile is a shapefile dataset type
	DatasetTypeShapefile

	// DatasetTypeNameGeopackage is the name of the geopackage dataset type
	DatasetTypeNameGeopackage DatasetTypeName = "geopackage"
	// DatasetTypeNameGeojson is the name of the geojson dataset type
	DatasetTypeNameGeojson DatasetTypeName = "geojson"
	// DatasetTypeNameShapefile is the name of the shapefile dataset type
	DatasetTypeNameShapefile DatasetTypeName = "shapefile"
)

func DatasetTypeNameToValue(typeName DatasetTypeName) DatasetTypeValue {
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

func DatasetTypeValueToName(typeValue DatasetTypeValue) DatasetTypeName {
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
