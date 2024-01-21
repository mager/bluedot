package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
	fs "github.com/mager/bluedot/firestore"
)

type StoreGeoJSONResp struct {
	Success bool `json:"success"`
}

func (h *Handler) storeGeoJSON(w http.ResponseWriter, r *http.Request) {
	resp := StoreGeoJSONResp{}
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

	// Fetch record from Firestore
	doc, err := h.Firestore.Collection("datasets").Doc(dataset.ID).Get(context.Background())
	if err != nil {
		h.Logger.Errorf("Error fetching Firestore record: %s", err)
		h.sendErrorJSON(w, http.StatusInternalServerError, "Error fetching Firestore record")
		return
	}

	ds := fs.Dataset{}
	err = doc.DataTo(&ds)
	if err != nil {
		h.Logger.Errorf("Error converting Firestore data to struct: %s", err)
		h.sendErrorJSON(w, http.StatusInternalServerError, "Error converting Firestore data to struct")
		return
	}

	if len(ds.Files) != 1 {
		h.Logger.Infof("Found %d files", len(ds.Files))
		return
	}

	file := ds.Files[0]
	geojsonBytes := getGeoJSONFromZipURLV3(file)
	h.Logger.Infof("GeoJSON: %s", geojsonBytes)

	resp.Success = true

}
