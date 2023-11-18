package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/everystreet/go-shapefile"
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

	// Print the filename and more details about the file
	fmt.Println(tmpfile.Name())
	fmt.Println(tmpfile.Stat())
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

	resp := GeoJSONResp{
		Type:     "FeatureCollection",
		Features: []Feature{},
	}

	for {
		record := scanner.Record()
		if record == nil {
			break
		}
		feature := record.GeoJSONFeature()
		fmt.Println(feature)

		// Set the geometry in the resp.Geometry map
		jsonData, err := json.Marshal(feature)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var f map[string]interface{}
		err = json.Unmarshal(jsonData, &f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the geometry in the resp.Geometry map
		resp.Features = append(resp.Features, Feature{
			Type:       "Feature",
			Properties: f["properties"].(map[string]interface{}),
			Geometry:   f["geometry"].(map[string]interface{}),
		})
	}

	// Err() returns the first error encountered during calls to Record()
	err = scanner.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: The GeoJSON files can be HUGE! Let's return a smaller response by simplifying the geometry
	// in each feature.

	// Set the Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
