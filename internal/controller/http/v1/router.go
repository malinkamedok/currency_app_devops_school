package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"p.solovev/internal/usecase"
)

func NewRouter(handler *chi.Mux, c usecase.CurrencyContract, i usecase.InfoContract) {
	handler.Route("/info", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Access-Control-Allow-Origin", "X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "origin", "x-requested-with"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
		NewInfoRoutes(r, c, i)
	})

	handler.Mount("/docs", http.StripPrefix("/docs", http.FileServer(http.Dir("./docs/swaggerui/"))))

	handler.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 - Not Found", http.StatusNotFound)
	})
}
