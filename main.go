package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"

	fs "cloud.google.com/go/firestore"
	gh "github.com/google/go-github/v56/github"
	"github.com/gorilla/mux"
	"github.com/mager/bluedot/config"
	"github.com/mager/bluedot/db"
	"github.com/mager/bluedot/firestore"
	"github.com/mager/bluedot/github"
	"github.com/mager/bluedot/handler"
	"github.com/mager/bluedot/logger"
	"github.com/mager/bluedot/router"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// @title Bluedot
// @version 1.0
// @description Primary backend for Geotory

// @contact.name @mager
// @contact.url https://geotory.com
// @contact.email magerleagues@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host api.geotory.com
// @BasePath /api
func main() {
	fx.New(
		fx.Provide(
			NewHTTPServer,
			zap.NewProduction,

			config.Options,
			db.Options,
			firestore.Options,
			github.Options,
			logger.Options,
			router.Options,
		),
		fx.Invoke(Register),
	).Run()

}

func Register(
	cfg config.Config,
	db *sql.DB,
	fs *fs.Client,
	gh *gh.Client,
	log *zap.SugaredLogger,
	router *mux.Router,
) {
	params := handler.Handler{
		Config:    cfg,
		Database:  db,
		Firestore: fs,
		Github:    gh,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
			},
		},
		Logger: log,
		Router: router,
	}

	handler.New(params)
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(handler.Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8085", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
