package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	fs "cloud.google.com/go/firestore"
	"github.com/k0kubun/pp"
	"github.com/mager/bluedot/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	geojson "github.com/paulmach/go.geojson"
)

type SaveFeaturesReq struct {
	URL     string `json:"url"`
	Dataset string `json:"dataset"`
}

type SaveFeaturesResp struct {
	ID      string                `json:"id"`
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Stats   SaveFeaturesRespStats `json:"stats"`
}

type SaveFeaturesRespStats struct {
	NumFeatures             int `json:"numFeatures"`
	NumPolygons             int `json:"numPolygons"`
	NumMultiPolygons        int `json:"numMultiPolygons"`
	NumFeaturesNotProcessed int `json:"numFeaturesNotProcessed"`
}

type SimplifyGeoJSONReq struct {
	Dataset string `json:"dataset"`
	URL     string `json:"url"`
}

const (
	tellusRespStatusSuccess = "success"
)

func (h *Handler) saveFeatures(w http.ResponseWriter, r *http.Request) {
	var req SaveFeaturesReq
	var resp SaveFeaturesResp
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := "https://tellus-zhokjvjava-uc.a.run.app/api/simplify/geojson"
	tellusReqBody := []byte(`{"url": "` + req.URL + `"}`)
	// Create a new HTTP request with the PUT method
	tellusResp, err := http.Post(p, "application/json", bytes.NewBuffer(tellusReqBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Close the response body
	defer tellusResp.Body.Close()

	body, err := ioutil.ReadAll(tellusResp.Body)
	if err != nil {
		panic(err)
	}

	var raw json.RawMessage
	err = json.Unmarshal(body, &raw)
	if err != nil {
		panic(err)
	}

	// Convert raw JSON to string
	jsonString := string(raw)
	fc, err := geojson.UnmarshalFeatureCollection([]byte(jsonString))
	if err != nil {
		panic(err)
	}

	// Fetch the dataset from Firestore
	dataset, err := h.Firestore.Collection("datasets").Doc(req.Dataset).Get(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	features := fc.Features

	numPolygons := 0
	numMultiPolygons := 0

	for _, f := range features {
		feature := firestore.Feature{}
		geos := []firestore.Geometry{}
		if f.Geometry.IsPolygon() {
			feature.Type = firestore.FeatureTypePolygon
			feat := f.Geometry.Polygon
			geos = ProcessPolygons(geos, feat)
			numPolygons++
		}

		if f.Geometry.IsMultiPolygon() {
			feature.Type = firestore.FeatureTypeMultiPolygon
			feat := f.Geometry.MultiPolygon
			for _, mp := range feat {
				geos = ProcessPolygons(geos, mp)
			}

			numMultiPolygons++
		}

		if f.Geometry.IsPoint() {
			pp.Print("TODO: Handle Point")
		}

		if f.Geometry.IsMultiPoint() {
			pp.Print("TODO: Handle MultiPoint")
		}

		if f.Geometry.IsLineString() {
			pp.Print("TODO: Handle LineString")
		}

		if f.Geometry.IsMultiLineString() {
			pp.Print("TODO: Handle MultiLineString")
		}

		// Look in the database to see if the signature exists
		u := firestore.GetFeatureUUID(f)

		// First look for the feature by UUID
		snap, err := h.Firestore.Collection("features").Doc(u.String()).Get(r.Context())
		if err != nil && status.Code(err) != codes.NotFound {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !snap.Exists() {
			feature.Dataset = req.Dataset
			feature.Geometries = geos
			feature.Name = GetPropertyName(f.Properties)
			feature.Properties = f.Properties

			_, err := h.Firestore.Collection("features").Doc(u.String()).Set(r.Context(), feature)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	numFeaturesNotProcessed := len(features) - numPolygons - numMultiPolygons
	resp.Stats = SaveFeaturesRespStats{
		NumFeatures:             len(features),
		NumPolygons:             numPolygons,
		NumMultiPolygons:        numMultiPolygons,
		NumFeaturesNotProcessed: numFeaturesNotProcessed,
	}

	// Update bounding box on the dataset
	_, err = h.Firestore.Collection("datasets").Doc(req.Dataset).Update(r.Context(), []fs.Update{
		{
			Path:  "bbox",
			Value: calculateBoundingBox(fc),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle response
	w.Header().Set("Content-Type", "application/json")
	resp.Success = true
	resp.ID = dataset.Ref.ID
	w.WriteHeader(http.StatusOK)

	// Return an error if there were unprocessed features
	if len(features) == 0 {
		resp.Message = "No features were processed"
		resp.Success = false
	}

	if numFeaturesNotProcessed != 0 {
		resp.Message = "Some features were not processed"
		resp.Success = false
	}

	// Write the response
	json.NewEncoder(w).Encode(resp)
}

func GetPropertyName(prop map[string]interface{}) string {
	var p string

	// List of possible keys for the property name
	keys := []string{"name", "NAME", "Name"}

	// Iterate through the keys and update the property name if found
	for _, key := range keys {
		if name, ok := prop[key]; ok {
			p = name.(string)
			break
		}
	}

	return p
}

func ProcessPolygons(geos []firestore.Geometry, feat [][][]float64) []firestore.Geometry {
	polyCoords := []float64{}
	for _, p := range feat {
		for _, c := range p {
			polyCoords = append(polyCoords, c...)
		}
	}
	geos = append(geos, firestore.Geometry{
		Coords: polyCoords,
	})

	return geos
}
