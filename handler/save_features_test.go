package handler

import (
	"reflect"
	"testing"

	"github.com/mager/bluedot/firestore"
)

func Test_ProcessPolygons(t *testing.T) {
	tests := []struct {
		name string
		geos []firestore.Geometry
		feat [][][]float64
		want []firestore.Geometry
	}{
		{
			name: "should add polygons to the geos array",
			geos: []firestore.Geometry{},
			feat: [][][]float64{
				{
					{1, 2},
					{3, 4},
					{5, 6},
				},
			},
			want: []firestore.Geometry{
				{
					Coords: []float64{1, 2, 3, 4, 5, 6},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessPolygons(tt.geos, tt.feat)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processPolygons() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetPropertyName(t *testing.T) {
	tests := []struct {
		name string
		prop map[string]interface{}
		want string
	}{
		{
			name: "should return the property name",
			prop: map[string]interface{}{
				"name": "test",
			},
			want: "test",
		},
		{
			name: "should return the property name",
			prop: map[string]interface{}{
				"NAME": "test",
			},
			want: "test",
		},
		{
			name: "should return the property name",
			prop: map[string]interface{}{
				"Name": "test",
			},
			want: "test",
		},
		{
			name: "should return an empty string",
			prop: map[string]interface{}{
				"test": "test",
			},
			want: "",
		},
		{
			name: "should return an empty string",
			prop: map[string]interface{}{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPropertyName(tt.prop); got != tt.want {
				t.Errorf("GetPropertyName() = %v, want %v", got, tt.want)
			}
		})
	}
}
