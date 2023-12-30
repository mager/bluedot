package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
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

func (h *Handler) saveFeatures(w http.ResponseWriter, r *http.Request) {
	var req SaveFeaturesReq
	var resp SaveFeaturesResp
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fc := getGeoJSONFromZipURL(req.URL)

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

	// Update bounding box and centroid on the dataset
	_, err = h.Firestore.Collection("datasets").Doc(req.Dataset).Update(r.Context(), []fs.Update{
		{
			Path:  "bbox",
			Value: h.calculateBoundingBox(fc),
		},
		{
			Path:  "centroid",
			Value: h.calculateCentroid(fc),
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

func (h *Handler) calculateBoundingBox(fc *geojson.FeatureCollection) [4]float64 {
	var bbox [4]float64

	coords := make([]float64, 0)
	for _, feature := range fc.Features {
		geom := feature.Geometry
		switch geom.Type {
		case geojson.GeometryPoint:
			coords = append(coords, geom.Point...)
		case geojson.GeometryPolygon:
			c := geom.Polygon[0]
			for _, coord := range c {
				coords = append(coords, coord...)
			}
		case geojson.GeometryMultiPolygon:
			c := geom.MultiPolygon[0][0]
			for _, coord := range c {
				coords = append(coords, coord...)
			}
		case geojson.GeometryLineString:
		case geojson.GeometryMultiPoint:
		case geojson.GeometryMultiLineString:
		case geojson.GeometryCollection:
			h.Logger.Info("TODO: Handle GeometryCollection")
		default:
			h.Logger.Info("Unknown geometry type")
		}
	}

	// Find the bounding box
	minX := coords[0]
	minY := coords[1]
	maxX := coords[0]
	maxY := coords[1]

	for i := 0; i < len(coords); i += 2 {
		x := coords[i]
		y := coords[i+1]

		if x < minX {
			minX = x
		}

		if y < minY {
			minY = y
		}

		if x > maxX {
			maxX = x
		}

		if y > maxY {
			maxY = y
		}

	}

	bbox[0] = minX
	bbox[1] = minY
	bbox[2] = maxX
	bbox[3] = maxY

	return bbox
}

func (h *Handler) calculateCentroid(fc *geojson.FeatureCollection) [2]float64 {
	var centroid [2]float64

	coords := make([]float64, 0)
	for _, feature := range fc.Features {
		geom := feature.Geometry
		switch geom.Type {
		case geojson.GeometryPoint:
			coords = append(coords, geom.Point...)
		case geojson.GeometryPolygon:
			c := geom.Polygon[0]
			for _, coord := range c {
				coords = append(coords, coord...)
			}
		case geojson.GeometryMultiPolygon:
			c := geom.MultiPolygon[0][0]
			for _, coord := range c {
				coords = append(coords, coord...)
			}
		case geojson.GeometryLineString:
		case geojson.GeometryMultiPoint:
		case geojson.GeometryMultiLineString:
		case geojson.GeometryCollection:
			log.Println("TODO: Handle GeometryCollection")
		default:
			log.Println("Unknown geometry type")
		}
	}

	var sumX float64
	var sumY float64
	for i := 0; i < len(coords); i += 2 {
		x := coords[i]
		y := coords[i+1]
		sumX += x
		sumY += y
	}

	centroid[0] = sumX / float64(len(coords)/2)
	centroid[1] = sumY / float64(len(coords)/2)

	return centroid
}

type PtolemyGeojsonReq struct {
	URL     string                `json:"url"`
	From    string                `json:"from"`
	Options PtolemyGeojsonOptions `json:"options"`
}

type PtolemyGeojsonOptions struct {
	Simplify PtolemyGeojsonOptionsSimplify `json:"simplify"`
}

type PtolemyGeojsonOptionsSimplify struct {
	Tolerance float64 `json:"tolerance"`
}

func getGeoJSONFromZipURL(url string) *geojson.FeatureCollection {
	ptolemyURL := "http://localhost:3005/api/geojson"
	// ptolemyURL := "https://ptolemy-zhokjvjava-uc.a.run.app/api/geojson"
	PtolemyGeojsonReqBody := PtolemyGeojsonReq{
		URL:  url,
		From: "shapefile",
		Options: PtolemyGeojsonOptions{
			Simplify: PtolemyGeojsonOptionsSimplify{
				Tolerance: 0.1,
			},
		},
	}

	// Convert the request body to JSON
	ptolemyReqBody, err := json.Marshal(PtolemyGeojsonReqBody)
	if err != nil {
		panic(err)
	}

	reqBody := bytes.NewBuffer(ptolemyReqBody)
	pp.Print(reqBody)
	ptolemyResp, err := http.Post(ptolemyURL, "application/json", reqBody)
	if err != nil {
		panic(err)
	}

	// Close the response body
	defer ptolemyResp.Body.Close()

	// If not 200, return an error
	if ptolemyResp.StatusCode != http.StatusOK {
		// Print the response body
		body, err := io.ReadAll(ptolemyResp.Body)
		if err != nil {
			panic(err)
		}
		panic(string(body))
	}

	body, err := io.ReadAll(ptolemyResp.Body)
	if err != nil {
		panic(err)
	}

	type PtolemyGeojsonRespErr struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var ptolemyGeojsonRespErr PtolemyGeojsonRespErr
	err = json.Unmarshal(body, &ptolemyGeojsonRespErr)
	if err != nil {
		panic(err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(body)
	if err != nil {
		panic(err)
	}
	return fc
}
