package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
)

// ServeHTTP handles an HTTP requests.
func (h *Handler) deleteDataset(w http.ResponseWriter, r *http.Request) {
	resp := DatasetResp{}
	vars := mux.Vars(r)
	username := vars["username"]
	datasetSlug := vars["slug"]

	user := db.GetUserByUsername(h.Database, username)
	if user.ID == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	dataset := db.GetDatasetByUserIdAndSlug(h.Database, user.ID, datasetSlug)
	if dataset.ID == "" {
		http.Error(w, "Dataset not found", http.StatusNotFound)
		return
	}

	// Delete dataset from database
	db.DeleteDatasetByUserIdAndSlug(h.Database, user.ID, datasetSlug)

	// Delete from Firestore
	_, err := h.Firestore.Collection("datasets").Doc(dataset.ID).Delete(r.Context())
	if err != nil {
		http.Error(w, "Error deleting dataset from Firestore", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
