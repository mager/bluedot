package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
)

type DeleteDatasetResp struct {
	ID string `json:"id"`
}

// deleteDataset godoc
//
//	@Summary		Delete a dataset
//	@Description	Delete a dataset
//	@ID				delete-dataset
//	@Tags			dataset
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string	true	"Username"
//	@Param			slug		path	string	true	"Slug"
//	@Success		200	{object}	DeleteDatasetResp
//	@Failure		404	{object}	ErrorResp
//	@Failure		500	{object}	ErrorResp
//	@Router			/datasets/{username}/{slug} [delete]
func (h *Handler) deleteDataset(w http.ResponseWriter, r *http.Request) {
	resp := DeleteDatasetResp{}
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

	// Delete dataset from database
	db.DeleteDatasetByUserIdAndSlug(h.Database, user.ID, datasetSlug)

	// Delete from Firestore
	_, err := h.Firestore.Collection("datasets").Doc(dataset.ID).Delete(r.Context())
	if err != nil {
		h.sendErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Delete features
	h.deleteFeaturesController(r.Context(), dataset.ID)

	h.Logger.Infow("Deleted dataset", "dataset", dataset.ID)

	resp.ID = dataset.ID

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
