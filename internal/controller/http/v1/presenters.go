package v1

import (
	"github.com/go-chi/chi/v5"
	"p.solovev/internal/usecase"
)

type infoRoutes struct {
	c usecase.CurrencyContract
	i usecase.InfoContract
}

func NewInfoRoutes(routes chi.Router, c usecase.CurrencyContract, i usecase.InfoContract) {
	ir := &infoRoutes{c: c, i: i}

	routes.Get("/", ir.getServiceInfo)
	routes.Get("/currency", ir.getCurrencyRate)
}
