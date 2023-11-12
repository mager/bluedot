package handler

import (
	"encoding/json"
	"net/http"
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

	resp.Type = "Feature"
	resp.Properties = make(map[string]interface{})
	resp.Properties["name"] = "dummy"
	resp.Properties["description"] = "dummy"
	resp.Geometry = make(map[string]interface{})
	resp.Geometry["type"] = "Point"
	resp.Geometry["coordinates"] = []float64{0, 0}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
