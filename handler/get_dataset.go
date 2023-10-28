package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mager/bluedot/db"
	"go.uber.org/zap"
)

// GetDataset is an http.Handler that copies its request body
// back to the response.
type GetDataset struct {
	log *zap.SugaredLogger
	sql *sql.DB
}

func (*GetDataset) Pattern() string {
	return "/datasets/"
}

// NewGetDataset builds a new GetDataset.
func NewGetDataset(log *zap.SugaredLogger, sql *sql.DB) *GetDataset {
	return &GetDataset{
		log: log,
		sql: sql,
	}
}

type GetDatasetResp struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Source      string `json:"source"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ServeHTTP handles an HTTP request to the /echo endpoint.
func (h *GetDataset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/datasets/")

	// TODO: Handle all routes manually

	// The URL will be /datasets/{username}/{datasetSlug}
	split := strings.Split(id, "/")
	username := split[0]
	datasetSlug := split[1]

	// Get the user from the database
	user := db.GetUserByUsername(h.sql, username)

	// Get the dataset from the database
	dataset := db.GetDatasetByUserIdAndSlug(h.sql, user.ID, datasetSlug)

	// Return the dataset object to the client
	var resp GetDatasetResp
	resp.ID = dataset.ID
	resp.UserID = dataset.UserID
	resp.Name = dataset.Name
	resp.Slug = dataset.Slug
	resp.Source = dataset.Source

	if dataset.Description.Valid {
		resp.Description = dataset.Description.String
	}
	if dataset.Created.Valid {
		// Log the dataset.Created.Time object to the console
		fmt.Println(dataset.Created.Time)
		resp.CreatedAt = dataset.Created.Time.Format("2006-01-02 15:04:05")
	}
	if dataset.Updated.Valid {
		resp.UpdatedAt = dataset.Updated.Time.Format("2006-01-02 15:04:05")
	}

	h.log.Infow("GetDataset", "dataset", resp)

	// Return in JSON format
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
}
