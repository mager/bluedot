package firestore

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	geojson "github.com/paulmach/go.geojson"
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

// Feature types
const (
	// FeatureTypeValueUnknown is an unknown feature type
	FeatureTypeValueUnknown int = iota
	// FeatureTypePoint is a point feature type
	FeatureTypePoint
	// FeatureTypeLineString is a line string feature type
	FeatureTypeLineString
	// FeatureTypePolygon is a polygon feature type
	FeatureTypePolygon
	// FeatureTypeMultiPoint is a multi point feature type
	FeatureTypeMultiPoint
	// FeatureTypeMultiLineString is a multi line string feature type
	FeatureTypeMultiLineString
	// FeatureTypeMultiPolygon is a multi polygon feature type
	FeatureTypeMultiPolygon
	// FeatureTypeGeometryCollection is a geometry collection feature type
	FeatureTypeGeometryCollection
)

type Dataset struct {
	Image    string    `json:"image" firestore:"image"`
	Source   string    `json:"source" firestore:"source"`
	Features []Feature `json:"features" firestore:"features"`
}

type Feature struct {
	Dataset    string      `json:"dataset" firestore:"dataset"`
	Type       int         `json:"type" firestore:"type"`
	Properties interface{} `json:"properties" firestore:"properties"`
	Geometries []Geometry  `json:"geometries" firestore:"geometries"`
	Name       string      `json:"name" firestore:"name"`
}

type Geometry struct {
	Coords []float64 `json:"coords" firestore:"coords"`
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

func GetFeatureUUID(f *geojson.Feature) uuid.UUID {
	// Create a map of the feature properties and the first coordinate of the first geometry
	data := map[string]interface{}{
		"properties": f.Properties,
		"firstCoord": GetFirstCoordinate(f),
	}

	// Serialize the map into a JSON string
	dataJSON, _ := json.Marshal(data)

	// Generate a UUID using SHA-1 hash of the namespace and serialized data
	u := uuid.NewSHA1(uuid.NameSpaceDNS, append([]byte("geotory"), dataJSON...))
	return u
}

func GetFirstCoordinate(f *geojson.Feature) []float64 {
	// Check feature type and return the first coordinate
	switch f.Geometry.Type {
	case geojson.GeometryPoint:
		return f.Geometry.Point
	case geojson.GeometryLineString:
		return f.Geometry.LineString[0]
	case geojson.GeometryPolygon:
		return f.Geometry.Polygon[0][0]
	case geojson.GeometryMultiPoint:
		return f.Geometry.MultiPoint[0]
	case geojson.GeometryMultiLineString:
		return f.Geometry.MultiLineString[0][0]
	case geojson.GeometryMultiPolygon:
		return f.Geometry.MultiPolygon[0][0][0]
	default:
		return f.BoundingBox
	}
}
