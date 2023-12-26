package firestore

import (
	"context"
	"log"
	"testing"

	"cloud.google.com/go/firestore"
	geojson "github.com/paulmach/go.geojson"
)

// ProvideFirestore provides a firestore client
func ProvideFirestoreBuilder() *firestore.Client {
	projectID := "geotory"

	client, err := firestore.NewClient(context.TODO(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func Test_GetFirstCoordinate(t *testing.T) {
	tests := []struct {
		name string
		feat *geojson.Feature
		want []float64
	}{
		{
			name: "should return the first coordinate for a geojson.GeometryPoint",
			feat: geojson.NewFeature(geojson.NewPointGeometry([]float64{1, 2})),
			want: []float64{1, 2},
		},
		{
			name: "should return the first coordinate for a geojson.GeometryLineString",
			feat: geojson.NewFeature(geojson.NewLineStringGeometry([][]float64{{1, 2}, {3, 4}})),
			want: []float64{1, 2},
		},
		{
			name: "should return the first coordinate for a geojson.GeometryPolygon",
			feat: geojson.NewFeature(
				geojson.NewPolygonGeometry(
					[][][]float64{
						{
							{1, 2},
							{3, 4},
							{5, 6},
						},
					},
				),
			),
			want: []float64{1, 2},
		},
		{
			name: "should return the first coordinate for a geojson.GeometryMultiPoint",
			feat: geojson.NewFeature(
				geojson.NewMultiPointGeometry(
					[]float64{1, 2},
					[]float64{3, 4},
				),
			),
			want: []float64{1, 2},
		},
		{
			name: "should return the first coordinate for a geojson.GeometryMultiLineString",
			feat: geojson.NewFeature(
				geojson.NewMultiLineStringGeometry(
					[][]float64{
						{1, 2},
						{3, 4},
					},
					[][]float64{
						{5, 6},
						{7, 8},
					},
				),
			),
			want: []float64{1, 2},
		},
		{
			name: "should return the first coordinate for a geojson.GeometryMultiPolygon",
			feat: geojson.NewFeature(
				geojson.NewMultiPolygonGeometry(
					[][][]float64{
						{
							{1, 2},
							{3, 4},
							{5, 6},
						},
					},
					[][][]float64{
						{
							{7, 8},
							{9, 10},
							{11, 12},
						},
					},
				),
			),
			want: []float64{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFirstCoordinate(tt.feat)

			if got[0] != tt.want[0] || got[1] != tt.want[1] {
				t.Errorf("GetFirstCoordinate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetFeatureUUID(t *testing.T) {
	tests := []struct {
		name string
		feat *geojson.Feature
		want string
	}{
		{
			name: "should return a UUID for a geojson.GeometryPoint",
			feat: geojson.NewFeature(geojson.NewPointGeometry([]float64{1, 2})),
			want: "32710693-87ba-5b48-b1d1-ed6b7308ba65",
		},
		{
			name: "should return a UUID for a geojson.GeometryLineString",
			feat: geojson.NewFeature(geojson.NewLineStringGeometry([][]float64{{1, 2}, {3, 4}})),
			want: "32710693-87ba-5b48-b1d1-ed6b7308ba65",
		},
		{
			name: "should return a UUID for a geojson.GeometryPolygon",
			feat: geojson.NewFeature(
				geojson.NewPolygonGeometry(
					[][][]float64{
						{
							{1, 2},
							{3, 4},
							{5, 6},
						},
					},
				),
			),
			want: "32710693-87ba-5b48-b1d1-ed6b7308ba65",
		},
		{
			name: "should return a UUID for a geojson.GeometryMultiPoint",
			feat: geojson.NewFeature(
				geojson.NewMultiPointGeometry(
					[]float64{1, 2},
					[]float64{3, 4},
				),
			),
			want: "32710693-87ba-5b48-b1d1-ed6b7308ba65",
		},
		{
			name: "should return a UUID for a geojson.GeometryMultiLineString",
			feat: geojson.NewFeature(
				geojson.NewMultiLineStringGeometry(
					[][]float64{
						{1, 2},
						{3, 4},
					},
					[][]float64{
						{5, 6},
						{7, 8},
					},
				),
			),
			want: "32710693-87ba-5b48-b1d1-ed6b7308ba65",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.feat.SetProperty("name", "test")
			got := GetFeatureUUID(tt.feat)

			if got.String() != tt.want {
				t.Errorf("GetFeatureUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DatasetTypeValueToName(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want string
	}{
		{
			name: "should return the name of the dataset type",
			arg:  1,
			want: "geopackage",
		},
		{
			name: "should return the name of the dataset type",
			arg:  2,
			want: "geojson",
		},
		{
			name: "should return the name of the dataset type",
			arg:  3,
			want: "shapefile",
		},
		{
			name: "should return an empty string",
			arg:  4,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DatasetTypeValueToName(tt.arg)

			if got != tt.want {
				t.Errorf("DatasetTypeValueToName() = %v, want %v", got, tt.want)
			}
		})
	}
}
