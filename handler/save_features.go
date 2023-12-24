package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

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

	for _, f := range features[10:13] {
		geos := []firestore.Geometry{}
		if f.Geometry.IsPolygon() {
			pp.Print("Found a polygon!")
			feat := f.Geometry.Polygon
			props := f.Properties
			pp.Print("Properties:", props)
			polyCoords := []float64{}
			for i, p := range feat {
				if i == 0 {
					for _, c := range p {
						polyCoords = append(polyCoords, c...)
					}
				} else {
					pp.Print("TODO: Handle holes")
				}
			}
			geos = append(geos, firestore.Geometry{
				Coords: polyCoords,
			})
			numPolygons++
			pp.Print("Coords for polygon:", polyCoords)
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

		// Look in the database to see if the signature exists
		u := firestore.GetFeatureUUID(f.Properties)

		// First look for the feature by UUID
		snap, err := h.Firestore.Collection("features").Doc(u.String()).Get(r.Context())
		if err != nil && status.Code(err) != codes.NotFound {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pp.Print("Found feature by UUID:", snap.Exists())

		// If the document already exists, then skip it
		if !snap.Exists() {
			// If the signature does not exist, then save the feature to Firestore
			feature := firestore.Feature{
				Dataset:    req.Dataset,
				Type:       firestore.FeatureTypePolygon,
				Properties: f.Properties,
				Geometries: geos,
			}

			_, err := h.Firestore.Collection("features").Doc(u.String()).Set(r.Context(), feature)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		pp.Println("Saved feature to Firestore!", f.Properties)
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
