package handler

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v56/github"
	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
)

// downloadDatasetZip godoc
//
//	@Summary		Download a dataset as a zip file
//	@Description	Download a dataset as a zip file
//	@ID				download-dataset-zip
//	@Tags			dataset
//	@Produce		application/zip
//	@Param			username	path	string	true	"Username"
//	@Param			slug		path	string	true	"Slug"
//	@Success		200
//	@Failure		404	{object}	ErrorResp
//	@Failure		500	{object}	ErrorResp
//	@Router			/datasets/{username}/{slug}/zip [get]
func (h *Handler) downloadDatasetZip(w http.ResponseWriter, r *http.Request) {
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

	// Call Github API to get the zip file
	owner, repo, path := parseGithubSource(dataset.Source)
	_, dc, _, err := h.Github.Repositories.GetContents(context.TODO(), owner, repo, path, &github.RepositoryContentGetOptions{
		Ref: "main",
	})
	if err != nil {
		h.Logger.Errorf("Error fetching contents: %s", err)
		h.sendErrorJSON(w, http.StatusInternalServerError, err.Error())
	}

	filename := fmt.Sprintf("%s.zip", datasetSlug)

	// Print the contents of the files returned
	for _, file := range dc {
		f := file.GetName()
		h.Logger.Infof("File found!: %s", f)
	}

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		h.Logger.Errorf("Error creating file: %s", err)
		h.sendErrorJSON(w, http.StatusInternalServerError, err.Error())
	}

	// Close the file
	defer out.Close()

	// Create a new zip writer
	zw := zip.NewWriter(out)
	defer zw.Close()

	// Iterate through the files and write them to the zip archive
	for _, file := range dc {
		fileData, ghResp, err := h.Github.Repositories.DownloadContents(context.TODO(), owner, repo, file.GetPath(), &github.RepositoryContentGetOptions{
			Ref: "main",
		})
		if err != nil {
			h.Logger.Errorf("Error fetching file contents: %s", err)
			continue
		}
		defer ghResp.Body.Close()

		// Create a zip file header
		fh := &zip.FileHeader{
			Name:   file.GetName(),
			Method: zip.Deflate,
		}

		// Write the file header and contents to the zip archive
		fw, err := zw.CreateHeader(fh)
		if err != nil {
			h.Logger.Errorf("Error creating zip file header: %s", err)
			continue
		}

		// Write the file contents to the zip archive
		_, err = io.Copy(fw, fileData)
		if err != nil {
			h.Logger.Errorf("Error writing file contents: %s", err)
			continue
		}
	}
	zw.Close()

	// Set the Content-Type and Content-Disposition headers
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	// Write the zip file to the response
	_, err = io.Copy(w, out)
	if err != nil {
		h.Logger.Errorf("Error writing file: %s", err)
		h.sendErrorJSON(w, http.StatusInternalServerError, err.Error())
	}

	// Return the file
	http.ServeFile(w, r, filename)
}
