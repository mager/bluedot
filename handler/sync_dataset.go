package handler

func (h *DatasetsHandler) syncDataset(resp *DatasetResp, username, datasetSlug string) *DatasetResp {
	h.log.Infof("Syncing dataset %s/%s", username, datasetSlug)

	resp = h.getDataset(resp, username, datasetSlug)

	// Fetch the filenames from the source

	return resp
}
