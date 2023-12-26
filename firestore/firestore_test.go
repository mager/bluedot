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
