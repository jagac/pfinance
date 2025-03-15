package routes

import (
	"net/http"

	"github.com/jagac/pfinance/internal/handlers"
)

type AssetRouter struct {
	handler *handlers.AssetHandler
}

func NewAssetRouter(handler *handlers.AssetHandler) *AssetRouter {
	return &AssetRouter{
		handler: handler,
	}
}

func (r *AssetRouter) RegisterRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.Handle("POST /api/assets", http.HandlerFunc(r.handler.CreateAsset))
	mux.Handle("GET /api/assets/{id}", http.HandlerFunc(r.handler.GetAsset))

	return mux
}
