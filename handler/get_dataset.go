package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mager/bluedot/db"
)

// ServeHTTP handles an HTTP requests.
func (h *DatasetsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If not /datasets/, return 404
	if !strings.HasPrefix(r.URL.Path, "/datasets/") {
		http.NotFound(w, r)
		return
	}

	// Get the dataset ID from the URL
	id := strings.TrimPrefix(r.URL.Path, "/datasets/")

	// Handle error cases
	if id == "" {
		http.Error(w, "No dataset ID provided", http.StatusBadRequest)
		return
	}

	// Handle when there isn't a slash in the ID
	if !strings.Contains(id, "/") {
		http.Error(w, "Invalid dataset ID provided", http.StatusBadRequest)
		return
	}

	// The URL will be /datasets/{username}/{datasetSlug}
	split := strings.Split(id, "/")
	username := split[0]
	datasetSlug := split[1]

	resp := DatasetResp{}

	if r.Method == http.MethodGet {
		// Get Datasets
		h.getDataset(&resp, username, datasetSlug)
	} else if r.Method == http.MethodPut {
		// Sync Datasets
		h.syncDataset(&resp, username, datasetSlug)
	}

	// Return in JSON format
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
}

func (h *DatasetsHandler) getDataset(resp *DatasetResp, username, datasetSlug string) *DatasetResp {
	// Get the user from the database
	user := db.GetUserByUsername(h.sql, username)

	// Get the dataset from the database
	dataset := db.GetDatasetByUserIdAndSlug(h.sql, user.ID, datasetSlug)

	// Set the response
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

	return resp
}
