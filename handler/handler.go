package handler

import (
	"database/sql"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/google/go-github/v56/github"
	"github.com/gorilla/mux"
	"github.com/mager/bluedot/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Route is an http.Handler that knows the mux pattern
// under which it will be registered.
type Route interface {
	http.Handler

	// Pattern reports the path at which this is registered.
	Pattern() string
}

// NewServeMux builds a ServeMux that will route requests
// to the given Route.
func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}

type Handler struct {
	fx.In

	Config    config.Config
	Database  *sql.DB
	Firestore *firestore.Client
	Github    *github.Client
	Logger    *zap.SugaredLogger
	Router    *mux.Router
}

// New creates a Handler struct
func New(h Handler) *Handler {
	h.registerRoutes()
	return &h
}

// RegisterRoutes registers all the routes for the route handler
func (h *Handler) registerRoutes() {
	h.Router.HandleFunc("/datasets/{username}", h.getDatasets).Methods("GET")
	h.Router.HandleFunc("/datasets/{username}/{slug}", h.getDataset).Methods("GET")
	h.Router.HandleFunc("/datasets/{username}/{slug}", h.syncDataset).Methods("PUT")
	h.Router.HandleFunc("/datasets/{username}/{slug}", h.deleteDataset).Methods("DELETE")
	h.Router.HandleFunc("/datasets/{username}/{slug}/zip", h.downloadDatasetZip).Methods("GET")

	// Experimental
	h.Router.HandleFunc("/shapefile", h.getShapefileJSON).Methods("GET")
}
