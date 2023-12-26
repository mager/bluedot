package handler

import (
	"database/sql"

	"cloud.google.com/go/firestore"
	"github.com/google/go-github/v56/github"
	geojson "github.com/paulmach/go.geojson"
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

	Image string        `json:"image"`
	Types []DatasetType `json:"types"`

	Geojson *geojson.FeatureCollection `json:"geojson"`
}

type DatasetType struct {
	Name string `json:"name"`
}

type Datasets struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Image string `json:"image"`
}

type DatasetsResp struct {
	Datasets []Datasets `json:"datasets"`
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
