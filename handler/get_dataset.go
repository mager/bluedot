package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// GetDataset is an http.Handler that copies its request body
// back to the response.
type GetDataset struct{}

// NewGetDataset builds a new GetDataset.
func NewGetDataset() *GetDataset {
	return &GetDataset{}
}

// ServeHTTP handles an HTTP request to the /echo endpoint.
func (*GetDataset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to handle request:", err)
	}
}
