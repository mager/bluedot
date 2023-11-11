package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
)

// ServeHTTP handles an HTTP requests.
func (h *Handler) getDataset(w http.ResponseWriter, r *http.Request) {
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

	resp.ID = dataset.ID
	resp.UserID = dataset.UserID
	resp.Name = dataset.Name
	resp.Slug = dataset.Slug
	resp.Source = dataset.Source

	if dataset.Description.Valid {
		resp.Description = dataset.Description.String
	}

	if dataset.Created.Valid {
		resp.CreatedAt = dataset.Created.Time.Format("2006-01-02 15:04:05")
	}

	if dataset.Updated.Valid {
		resp.UpdatedAt = dataset.Updated.Time.Format("2006-01-02 15:04:05")
	}

	resp.User.Image = user.Image
	resp.User.Slug = user.Slug

	// Fetch record from Firestore
	doc, err := h.Firestore.Collection("datasets").Doc(resp.ID).Get(context.Background())
	if err != nil {
		h.Logger.Errorf("Error fetching Firestore record: %s", err)
		http.Error(w, "Error fetching Firestore record", http.StatusInternalServerError)
		return
	}

	resp.Image = doc.Data()["image"].(string)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
