package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	fs "cloud.google.com/go/firestore"

	"github.com/google/go-github/v56/github"
	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
	"github.com/mager/bluedot/firestore"
)

// syncDataset godoc
//
//	@Summary		Sync a dataset
//	@Description	Syncing a dataset
//	@ID				sync-dataset
//	@Tags			dataset
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string	true	"Username"
//	@Param			slug		path	string	true	"Slug"
//	@Success		200	{object}	DatasetResp
//	@Failure		404	{object}	ErrorResp
//	@Failure		500	{object}	ErrorResp
//	@Router			/datasets/{username}/{slug} [put]
func (h *Handler) syncDataset(w http.ResponseWriter, r *http.Request) {
	resp := DatasetResp{}
	vars := mux.Vars(r)
	username := vars["username"]
	datasetSlug := vars["slug"]

	user := db.GetUserByUsername(h.Database, username)
	if user.ID == "" {
		h.sendErrorJSON(w, http.StatusNotFound, "User not found")
		return
	}

	dataset := db.GetDatasetByUserIdAndSlug(h.Database, user.ID, datasetSlug)
	if dataset.ID == "" {
		h.sendErrorJSON(w, http.StatusNotFound, "Dataset not found")
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
	files := []string{}
	record := map[string]interface{}{
		"source": fmt.Sprintf("%s/%s/%s", owner, repo, path),
		"types":  types,
	}

	zips := []string{}
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
			record["types"] = append(record["types"].([]int), firestore.DatasetTypeGeopackage)
		}
		if f[len(f)-8:] == ".geojson" {
			record["types"] = append(record["types"].([]int), firestore.DatasetTypeGeojson)
		}

		if f[len(f)-4:] == ".zip" {
			files = append(files, file.GetDownloadURL())
		}
	}
	record["files"] = files

	h.Logger.Infow("Record", "record", record, "zips", zips)

	// If there is a zip file, we need to parse through it and save features
	if len(files) == 1 {
		for _, file := range files {
			if strings.Contains(file, ".zip") {
				fc := getGeoJSONFromZipURL(file)
				record["bbox"] = h.calculateBoundingBox(fc)
				record["centroid"] = h.calculateCentroid(fc)
			}
		}
	}
	// 	fc := getGeoJSONFromZipURL(zips[0])
	// 	features := fc.Features

	// 	numPolygons := 0
	// 	numMultiPolygons := 0

	// 	for _, f := range features {
	// 		feature := firestore.Feature{}
	// 		geos := []firestore.Geometry{}
	// 		if f.Geometry.IsPolygon() {
	// 			feature.Type = firestore.FeatureTypePolygon
	// 			feat := f.Geometry.Polygon
	// 			geos = ProcessPolygons(geos, feat)
	// 			numPolygons++
	// 		}

	// 		if f.Geometry.IsMultiPolygon() {
	// 			feature.Type = firestore.FeatureTypeMultiPolygon
	// 			feat := f.Geometry.MultiPolygon
	// 			for _, mp := range feat {
	// 				geos = ProcessPolygons(geos, mp)
	// 			}

	// 			numMultiPolygons++
	// 		}

	// 		if f.Geometry.IsPoint() {
	// 			pp.Print("TODO: Handle Point")
	// 		}

	// 		if f.Geometry.IsMultiPoint() {
	// 			pp.Print("TODO: Handle MultiPoint")
	// 		}

	// 		if f.Geometry.IsLineString() {
	// 			pp.Print("TODO: Handle LineString")
	// 		}

	// 		if f.Geometry.IsMultiLineString() {
	// 			pp.Print("TODO: Handle MultiLineString")
	// 		}

	// 		// Look in the database to see if the signature exists
	// 		u := firestore.GetFeatureUUID(f)

	// 		// First look for the feature by UUID
	// 		snap, err := h.Firestore.Collection("features").Doc(u.String()).Get(r.Context())
	// 		if err != nil && status.Code(err) != codes.NotFound {
	// 			h.sendErrorJSON(w, http.StatusInternalServerError, "Error fetching Firestore record")
	// 			return
	// 		}

	// 		if !snap.Exists() {
	// 			feature.Dataset = dataset.ID
	// 			feature.Geometries = geos
	// 			feature.Name = GetPropertyName(f.Properties)
	// 			feature.Properties = f.Properties

	// 			_, err := h.Firestore.Collection("features").Doc(u.String()).Set(r.Context(), feature)
	// 			if err != nil {
	// 				h.sendErrorJSON(w, http.StatusInternalServerError, err.Error())
	// 				return
	// 			}
	// 		}
	// 	}

	// 	numFeaturesNotProcessed := len(features) - numPolygons - numMultiPolygons
	// 	h.Logger.Infow("Features", "numPolygons", numPolygons, "numMultiPolygons", numMultiPolygons, "numFeaturesNotProcessed", numFeaturesNotProcessed)

	// 	record["bbox"] = h.calculateBoundingBox(fc)
	// 	record["centroid"] = h.calculateCentroid(fc)
	// }

	// Create or update a record in Firestore
	h.Logger.Infof("Dataset ID: %s", dataset.ID)
	_, err = h.Firestore.Collection("datasets").
		Doc(dataset.ID).
		Set(context.Background(), record, fs.MergeAll)

	if err != nil {
		h.Logger.Errorf("Error updating Firestore: %s", err)
		h.sendErrorJSON(w, http.StatusInternalServerError, "Error updating Firestore")
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
