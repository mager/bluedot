package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"

	"github.com/mager/bluedot/config"
	"github.com/mager/bluedot/db"
	"github.com/mager/bluedot/github"
	"github.com/mager/bluedot/handler"
	"github.com/mager/bluedot/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			NewHTTPServer, handler.NewServeMux, zap.NewProduction,
			config.Options, db.Options, github.Options, logger.Options,

			// Handlers
			fx.Annotate(
				handler.NewGetDataset,
				fx.As(new(handler.Route)),
			),
		),
		fx.Invoke(func(*http.Server, config.Config, *sql.DB, *zap.SugaredLogger) {}, func() {
			fmt.Println("Hello, world!")
		}),
	).Run()

}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
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
