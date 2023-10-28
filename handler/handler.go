package handler

import "net/http"

// NewServeMux builds a ServeMux that will route requests
// to the given GetDataset.
func NewServeMux(getDataset *GetDataset) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/dataset", getDataset)
	return mux
}
