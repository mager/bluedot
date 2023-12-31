package handler

import (
	"context"
	"encoding/json"
	"net/http"
)

type DeleteFeaturesReq struct {
	Dataset string `json:"dataset"`
}

type DeleteFeaturesResp struct {
	ID string `json:"id"`
}

// deleteFeatures godoc
//
//	@Summary		Delete features
//	@Description	Delete features for a given dataset
//	@ID				delete-features
//	@Tags			dataset
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string				true	"Username"
//	@Param			slug		path	string				true	"Slug"
//	@Param 			request 	body 	DeleteFeaturesReq 	true 	"Delete Features Req"
//	@Success		200	{object}	DatasetResp
//	@Failure		404	{object}	ErrorResp
//	@Failure		500	{object}	ErrorResp
//	@Router			/datasets/{username}/{slug}/deleteFeatures [post]
func (h *Handler) deleteFeatures(w http.ResponseWriter, r *http.Request) {
	var req DeleteFeaturesReq
	var resp DeleteFeaturesResp
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deletedFeatures := h.deleteFeaturesController(r.Context(), req.Dataset)

	h.Logger.Infow("Deleted features", "dataset", req.Dataset, "features", deletedFeatures)
	resp.ID = req.Dataset

	w.Header().Set("Content-Type", "application/json")

	// Write the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) deleteFeaturesController(ctx context.Context, dataset string) []string {
	var deletedFeatures []string
	iter := h.Firestore.Collection("features").Where("dataset", "==", dataset).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		deletedFeatures = append(deletedFeatures, doc.Ref.ID)
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			h.Logger.Errorf("Error deleting feature: %s", err)
		}
	}
	return deletedFeatures
}
