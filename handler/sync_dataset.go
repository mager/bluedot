package handler

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/google/go-github/v56/github"
)

func (h *DatasetsHandler) syncDataset(resp *DatasetResp, username, datasetSlug string) *DatasetResp {
	h.log.Infof("Syncing dataset %s/%s", username, datasetSlug)

	resp = h.getDataset(resp, username, datasetSlug)

	// Fetch the filenames from the source
	fc, dc, ghr, err := h.gh.Repositories.GetContents(context.TODO(), "mager", "maps", "illinois", &github.RepositoryContentGetOptions{
		Ref: "main",
	})

	if err != nil {
		h.log.Errorf("Error fetching contents: %s", err)

	}

	h.log.Infof("Contents: %s", fc)
	h.log.Infof("Directory Contents: %s", dc)
	h.log.Infof("Github Response: %s", ghr)

	pngFile := ""

	// If there is a png file, extract the download URL
	for _, file := range dc {
		if file.GetName() == datasetSlug+".png" {
			pngFile = file.GetDownloadURL()
		}
	}

	h.log.Infof("PNG File found!: %s", pngFile)

	// Update a record in Firestore
	_, err = h.fs.Collection("datasets").Doc(resp.ID).Set(context.TODO(), map[string]interface{}{
		"image": pngFile,
	}, firestore.MergeAll)

	if err != nil {
		h.log.Errorf("Error updating Firestore: %s", err)
	}

	return resp
}
