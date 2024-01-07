package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	fs "cloud.google.com/go/firestore"

	"github.com/google/go-github/v56/github"
	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
	"github.com/mager/bluedot/firestore"
	geojson "github.com/paulmach/go.geojson"
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

type PtolemyGeojsonReq struct {
	URL     string                `json:"url"`
	From    string                `json:"from"`
	Options PtolemyGeojsonOptions `json:"options"`
}

type PtolemyGeojsonResp struct {
	GeoJSON *geojson.FeatureCollection `json:"geojson"`
	Context PtolemyGeojsonRespContext  `json:"context"`
}

type PtolemyGeojsonRespContext struct {
	Centroid []float64 `json:"centroid"`
	Zoom     int       `json:"zoom"`
}

type PtolemyGeojsonOptions struct {
	Simplify PtolemyGeojsonOptionsSimplify `json:"simplify"`
}

type PtolemyGeojsonOptionsSimplify struct {
	Tolerance float64 `json:"tolerance"`
}

func getGeoJSONFromZipURLV2(url string) *PtolemyGeojsonResp {
	ptolemyURL := "https://ptolemy-zhokjvjava-uc.a.run.app/api/geojson"
	PtolemyGeojsonReqBody := PtolemyGeojsonReq{
		URL:  url,
		From: "shapefile",
		Options: PtolemyGeojsonOptions{
			Simplify: PtolemyGeojsonOptionsSimplify{
				Tolerance: 0.01,
			},
		},
	}

	// Convert the request body to JSON
	ptolemyReqBody, err := json.Marshal(PtolemyGeojsonReqBody)
	if err != nil {
		panic(err)
	}

	reqBody := bytes.NewBuffer(ptolemyReqBody)
	ptolemyResp, err := http.Post(ptolemyURL, "application/json", reqBody)
	if err != nil {
		panic(err)
	}

	// Close the response body
	defer ptolemyResp.Body.Close()

	// If not 200, return an error
	if ptolemyResp.StatusCode != http.StatusOK {
		// Print the response body
		body, err := io.ReadAll(ptolemyResp.Body)
		if err != nil {
			panic(err)
		}
		panic(string(body))
	}

	body, err := io.ReadAll(ptolemyResp.Body)
	if err != nil {
		panic(err)
	}

	type PtolemyGeojsonRespErr struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var ptolemyGeojsonRespErr PtolemyGeojsonRespErr
	err = json.Unmarshal(body, &ptolemyGeojsonRespErr)
	if err != nil {
		panic(err)
	}

	var ptolemyGeojsonResp PtolemyGeojsonResp
	err = json.Unmarshal(body, &ptolemyGeojsonResp)
	if err != nil {
		panic(err)
	}

	return &ptolemyGeojsonResp
}
