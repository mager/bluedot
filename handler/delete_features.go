package handler

import (
	"encoding/json"
	"net/http"

	"github.com/k0kubun/pp"
)

type DeleteFeaturesReq struct {
	Dataset string `json:"dataset"`
}

type DeleteFeaturesResp struct {
	ID string `json:"id"`
}

func (h *Handler) deleteFeatures(w http.ResponseWriter, r *http.Request) {
	var req DeleteFeaturesReq
	var resp DeleteFeaturesResp
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var deletedFeatures []string
	iter := h.Firestore.Collection("features").Where("dataset", "==", req.Dataset).Documents(r.Context())
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		deletedFeatures = append(deletedFeatures, doc.Ref.ID)
		_, err = doc.Ref.Delete(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	pp.Print("Deleted features", "features", deletedFeatures)
	resp.ID = "success"

	w.Header().Set("Content-Type", "application/json")

	// Write the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
