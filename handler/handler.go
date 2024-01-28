package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
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

	Config     config.Config
	Database   *sql.DB
	Firestore  *firestore.Client
	Github     *github.Client
	HttpClient *http.Client
	Logger     *zap.SugaredLogger
	Router     *mux.Router
	Storage    *storage.Client
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
	h.Router.HandleFunc("/datasets/{username}/{slug}/deleteFeatures", h.deleteFeatures).Methods("POST")

	h.Router.HandleFunc("/datasets/{username}/{slug}/geojson", h.storeGeoJSON).Methods("POST")
	// Deprecated
}

type ErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *Handler) sendErrorJSON(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResp{
		Code:    code,
		Message: message,
	})
}
