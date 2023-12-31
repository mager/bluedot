package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
)

// getDatasets godoc
//
//	@Summary		Get all datasets for a user
//	@Description	Fetch datasets from a given user
//	@ID				get-datasets
//	@Tags			dataset
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string	true	"Username"
//	@Success		200	{object}	DatasetsResp
//	@Failure		404	{object}	ErrorResp
//	@Failure		500	{object}	ErrorResp
//	@Router			/datasets/{username} [get]
func (h *Handler) getDatasets(w http.ResponseWriter, r *http.Request) {
	resp := DatasetsResp{}
	vars := mux.Vars(r)
	username := vars["username"]

	user := db.GetUserByUsername(h.Database, username)
	if user.ID == "" {
		h.sendErrorJSON(w, http.StatusNotFound, "User not found")
		return
	}

	datasets := db.GetDatasetsByUserId(h.Database, user.ID)
	if len(datasets) == 0 {
		h.sendErrorJSON(w, http.StatusNotFound, "No datasets found")
		return
	}

	// Adapt the datasets to the response format
	for _, dataset := range datasets {
		datasetsResp := Datasets{
			ID:   dataset.ID,
			Name: dataset.Name,
			Slug: dataset.Slug,
		}
		resp.Datasets = append(resp.Datasets, datasetsResp)
	}

	// Fetch documents from Firestore
	datasetIDs := []string{}
	for _, dataset := range datasets {
		datasetIDs = append(datasetIDs, dataset.ID)
	}

	var documents []*firestore.DocumentSnapshot
	for _, datasetID := range datasetIDs {
		docRef := h.Firestore.Collection("datasets").Doc(datasetID)
		docSnap, err := docRef.Get(context.Background())
		if err != nil {
			h.Logger.Errorf("Error fetching Firestore record: %s", err)
			h.sendErrorJSON(w, http.StatussnringalServerError, "Error fetching Firestore record")
			return
		}
		documents = append(documents, docSnap)
	}

	// Find the images from the documents, and append to resp
	for _, doc := range documents {
		image := doc.Data()["image"].(string)
		for i, dataset := range resp.Datasets {
			if dataset.ID == doc.Ref.ID {
				resp.Datasets[i].Image = image
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
