package routes

import (
	"net/http"

	"github.com/jagac/pfinance/internal/handlers"
)

type AssetRouter struct {
	handler *handlers.AssetHandler
	logMiddleware func(http.Handler) http.Handler
	corsMiddleware func(http.Handler) http.Handler
}

func NewAssetRouter(handler *handlers.AssetHandler, logMiddleware func(http.Handler) http.Handler, corsMiddleware func(http.Handler) http.Handler) *AssetRouter {
	return &AssetRouter{
		handler: handler,
		logMiddleware: logMiddleware,
		corsMiddleware: corsMiddleware,
	}
}

func (r *AssetRouter) RegisterRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.Handle("POST /api/assets/new", r.corsMiddleware(r.logMiddleware(http.HandlerFunc(r.handler.CreateAsset))))
	mux.Handle("GET /api/assets/all", r.corsMiddleware(r.logMiddleware(http.HandlerFunc(r.handler.GetAssets))))
	mux.Handle("GET /api/returns", r.corsMiddleware(r.logMiddleware(http.HandlerFunc(r.handler.GetReturns))))
	mux.Handle("GET /api/returns/month", r.corsMiddleware(r.logMiddleware(http.HandlerFunc(r.handler.GetMonthlyReturns))))
	return mux
}
