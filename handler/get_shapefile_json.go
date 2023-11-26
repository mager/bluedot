package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/everystreet/go-shapefile"
	"github.com/paulmach/orb/geojson"
)

type GeoJSONResp struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Geometry   map[string]interface{} `json:"geometry"`
}

// ServeHTTP handles an HTTP requests.
func (h *Handler) getShapefileJSON(w http.ResponseWriter, r *http.Request) {
	// Fetch the Shapefile and then convert it to GeoJSON
	url := "https://github.com/mager/maps/raw/main/illinois-elections/tl_2012_17_vtd10.zip"
	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error fetching Shapefile", http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	// Create a temporary file to store the downloaded shapefile
	tmpfile, err := os.CreateTemp("", "shpp")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tmpfile.Close()

	_, err = io.Copy(tmpfile, response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.Open(tmpfile.Name())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scanner, err := shapefile.NewZipScanner(file, stat.Size(), "tl_2012_17_vtd10.zip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = scanner.Scan()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fc := geojson.NewFeatureCollection()
	counter := 0
	for counter < 3 {
		record := scanner.Record()
		if record == nil {
			break
		}
		feature := record.GeoJSONFeature()

		// Set the geometry in the resp.Geometry map
		jsonData, err := json.Marshal(feature)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var feat geojson.Feature
		err = json.Unmarshal(jsonData, &feat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fc.Append(&feat)
		counter++
	}

	// Err() returns the first error encountered during calls to Record()
	err = scanner.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fc)

}
