package v1

import (
	"log"
	"net/http"

	"go.uber.org/zap"
	"p.solovev/pkg/logger"

	"p.solovev/pkg/web"

	"github.com/go-chi/render"
)

type respCurrency struct {
	Data    map[string]float64 `json:"data"`
	Service string             `json:"service"`
}

func (ir *infoRoutes) getServiceInfo(w http.ResponseWriter, r *http.Request) {
	logger.Info("Got GET /info request")
	render.JSON(w, r, ir.i.GetInfoAboutService())
}

func (ir *infoRoutes) getCurrencyRate(w http.ResponseWriter, r *http.Request) {
	currencyCode := r.URL.Query().Get("currency")
	date := r.URL.Query().Get("date")

	logger.Info("Got GET /info/currency request with ", zap.String("currency", currencyCode), zap.String("date", date))

	response, err := ir.c.GetCurrencyRate(currencyCode, date)
	if err != nil {
		err := render.Render(w, r, web.ErrRender(err))
		if err != nil {
			log.Printf("Rendering error")
			return
		}
		return
	}

	responseJSON := respCurrency{Data: response, Service: "currency"}
	render.JSON(w, r, responseJSON)
}
