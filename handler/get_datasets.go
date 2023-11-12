package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
)

// ServeHTTP handles an HTTP requests.
func (h *Handler) getDatasets(w http.ResponseWriter, r *http.Request) {
	resp := DatasetsResp{}
	vars := mux.Vars(r)
	username := vars["username"]

	user := db.GetUserByUsername(h.Database, username)
	if user.ID == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	datasets := db.GetDatasetsByUserId(h.Database, user.ID)
	if len(datasets) == 0 {
		http.Error(w, "No datasets found", http.StatusNotFound)
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
			http.Error(w, "Error fetching Firestore record", http.StatusInternalServerError)
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
