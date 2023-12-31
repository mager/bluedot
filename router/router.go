package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProvideRouter provides a gorilla mux router
func ProvideRouter(lc fx.Lifecycle, log *zap.SugaredLogger) *mux.Router {
	var router = mux.NewRouter()

	router.Use(
		jsonMiddleware,
		// loggingDecorator(log),
	)

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				addr := ":8085"
				log.Info("Listening on ", addr)

				go http.ListenAndServe(addr, router)

				return nil
			},
		},
	)

	return router
}

// jsonMiddleware makes sure that every response is JSON
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// loggingDecorator logs the request
func loggingDecorator(log *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info(r.Method, " ", r.URL)
			next.ServeHTTP(w, r)
		})
	}
}

var Options = ProvideRouter
