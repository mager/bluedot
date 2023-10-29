package handler

import (
	"database/sql"

	"cloud.google.com/go/firestore"
	"github.com/google/go-github/v56/github"
	"go.uber.org/zap"
)

type DatasetResp struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Source      string `json:"source"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`

	User struct {
		Image string `json:"image"`
		Slug  string `json:"slug"`
	} `json:"user"`
}

// DatasetsHandler is an http.DatasetsHandler that copies its request body
// back to the response.
type DatasetsHandler struct {
	log *zap.SugaredLogger
	sql *sql.DB
	gh  *github.Client
	fs  *firestore.Client
}

func (*DatasetsHandler) Pattern() string {
	return "/datasets/"
}

// NewDatasetsHandler builds a new GetDataset.
func NewDatasetsHandler(log *zap.SugaredLogger, sql *sql.DB, gh *github.Client, fs *firestore.Client) *DatasetsHandler {
	return &DatasetsHandler{
		log: log,
		sql: sql,
		gh:  gh,
		fs:  fs,
	}
}
