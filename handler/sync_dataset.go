package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/google/go-github/v56/github"
	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
	fs "github.com/mager/bluedot/firestore"
)

func (h *Handler) syncDataset(w http.ResponseWriter, r *http.Request) {
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

	// Extract the owner, repo, and path from the Github source
	owner, repo, path := parseGithubSource(dataset.Source)

	// Fetch the filenames from the source
	_, dc, _, err := h.Github.Repositories.GetContents(context.TODO(), owner, repo, path, &github.RepositoryContentGetOptions{
		Ref: "main",
	})

	if err != nil {
		h.Logger.Errorf("Error fetching contents: %s", err)
	}

	var types []int
	record := map[string]interface{}{
		"source": fmt.Sprintf("%s/%s/%s", owner, repo, path),
		"types":  types,
	}

	for _, file := range dc {
		f := file.GetName()
		// If there is a filename ending in .png, set it as the image
		if f[len(f)-4:] == ".png" {
			record["image"] = f
		}
		// Use svg as backup
		if f[len(f)-4:] == ".svg" {
			record["image"] = f
		}
		// Handle types
		if f[len(f)-5:] == ".gpkg" {
			record["types"] = append(record["types"].([]int), fs.DatasetTypeGeopackage)
		}
		if f[len(f)-8:] == ".geojson" {
			record["types"] = append(record["types"].([]int), fs.DatasetTypeGeojson)
		}
	}

	// Create or update a record in Firestore
	h.Logger.Infof("Dataset ID: %s", dataset.ID)
	_, err = h.Firestore.Collection("datasets").
		Doc(dataset.ID).
		Set(context.Background(), record, firestore.MergeAll)

	if err != nil {
		h.Logger.Errorf("Error updating Firestore: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func parseGithubSource(source string) (string, string, string) {
	// Split the source into owner, repo, and path
	// Example input: mager/maps/illinois
	// Example output: mager, maps, illinois
	owner := strings.Split(source, "/")[0]
	repo := strings.Split(source, "/")[1]
	path := strings.Split(source, "/")[2]

	return owner, repo, path
}