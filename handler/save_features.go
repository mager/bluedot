package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/k0kubun/pp"
	geojson "github.com/paulmach/go.geojson"
)

type SaveFeaturesReq struct {
	URL     string `json:"url"`
	Dataset string `json:"dataset"`
}

type SaveFeaturesResp struct {
	ID string `json:"id"`
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

	pp.Println(dataset.Data())

	features := fc.Features
	pp.Print("Found ", len(features), " features")

	numPolygons := 0
	numMultiPolygons := 0

	for _, f := range features {
		if f.Geometry.IsPolygon() {
			pp.Print("Found a polygon!")
			polygon := f.Geometry.Polygon
			props := f.Properties
			pp.Print("Properties:", props)
			rawCoords := make([]float64, len(p)*2)
			for i, p := range polygon {
				pp.Print("Printing coordinates for poly index:", i)
				// For each polygon, we will have a long list of lat/long coordinates
				// but it will be as a single []float64
				for _, c := range p {
					rawCoords = append(rawCoords, c...)
				}
			}
			numPolygons++
			pp.Print("Raw coords for polygon:", rawCoords)
		}

		if f.Geometry.IsMultiPolygon() {
			pp.Print("TODO: Handle MultiPolygon")
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
	}

	pp.Println("Number of polygons:", numPolygons)
	pp.Println("Number of multi-polygons:", numMultiPolygons)
	pp.Println("Number of features:", len(features))

	resp.ID = "success"
	// // Get the properties
	// props := f.Properties

	// // Get the geometry
	// geom := f.Geometry

	// // Print the properties and geometry
	// fmt.Printf("properties: %v\n", props)
	// fmt.Printf("geometry: %v\n", geom)

	w.Header().Set("Content-Type", "application/json")

	// Write the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func getPropertyName(prop map[string]interface{}) string {
	p := ""

	if name, ok := prop["name"]; ok {
		p = name.(string)
	}

	if name, ok := prop["NAME"]; ok {
		p = name.(string)
	}

	if name, ok := prop["Name"]; ok {
		p = name.(string)
	}

	return p
}
