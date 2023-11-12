package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/jonas-p/go-shp"
)

// GeneralGeoJSONStruct is a generic GeoJSON struct
type GeneralGeoJSONStruct struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Geometry   map[string]interface{} `json:"geometry"`
}

// ServeHTTP handles an HTTP requests.
func (h *Handler) getShapefileJSON(w http.ResponseWriter, r *http.Request) {
	resp := GeneralGeoJSONStruct{}

	// Insert some dummy data
	// resp.Type = "Feature"
	// resp.Properties = make(map[string]interface{})
	// resp.Properties["name"] = "dummy"
	// resp.Properties["description"] = "dummy"
	// resp.Geometry = make(map[string]interface{})
	// resp.Geometry["type"] = "Point"
	// resp.Geometry["coordinates"] = []float64{0, 0}

	// Fetch the Shapefile and then convert it to GeoJSON
	url := "https://github.com/mager/maps/raw/main/illinois-elections/tl_2012_17_vtd10.shp"
	// Make an HTTP request to the URL
	// Convert the Shapefile to GeoJSON
	// Return the GeoJSON

	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error fetching Shapefile", http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	// Create a temporary file to store the downloaded data
	tempFile, err := os.CreateTemp("", "shapefile-*")
	if err != nil {
		h.Logger.Errorf("Error creating temporary file: %s", err)
	}
	defer os.Remove(tempFile.Name())

	// Write the data to the temporary file
	_, err = io.Copy(tempFile, response.Body)
	if err != nil {
		h.Logger.Errorf("Error writing to temporary file: %s", err)
	}

	// Parse the downloaded shapefile
	shape, err := shp.Open(tempFile.Name())
	if err != nil {
		h.Logger.Errorf("Error opening shapefile: %s", err)
	}

	// Get the fields from the shapefile
	fields := shape.Fields()

	// Get the shapefile records
	for shape.Next() {
		n, p := shape.Shape()
		h.Logger.Infof("Shape: %d, %v", n, p)

		for k, f := range fields {
			val := shape.ReadAttribute(n, k)
			h.Logger.Infof("Field: %s, %v", f, val)
		}
	}

	// Set the Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
