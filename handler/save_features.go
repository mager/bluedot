package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SaveFeaturesReq struct {
	URL string `json:"url"`
}

type SaveFeaturesResp struct {
	ID string `json:"id"`
}

type SimplifyGeoJSONReq struct {
	URL string `json:"url"`
}

type SimplifyGeoJSONResp struct {
	Status string      `json:"status"`
	Data   GeoJSONResp `json:"data"`
}

func (h *Handler) saveFeatures(w http.ResponseWriter, r *http.Request) {
	resp := SaveFeaturesResp{}

	var req SaveFeaturesReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := "https://tellus-zhokjvjava-uc.a.run.app/api/simpliify/geojson"
	tellusReqBody := []byte(`{"url": "` + req.URL + `"}`)
	// Create a new HTTP request with the PUT method
	tellusResp, err := http.Post(p, "application/json", bytes.NewBuffer(tellusReqBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode the JSON response body into GeoJSONResp
	var tellusRespBody SimplifyGeoJSONResp
	err = json.NewDecoder(tellusResp.Body).Decode(&tellusRespBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("tellusRespBody: %v", tellusRespBody.Status)

	// // Save the features to the database
	// featureID, _, err := h.Firestore.Collection("features").Add(r.Context(), map[string]interface{}{
	// 	"type": 2,
	// })

	// resp.ID = featureID.ID

	w.Header().Set("Content-Type", "application/json")

	// Write the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
