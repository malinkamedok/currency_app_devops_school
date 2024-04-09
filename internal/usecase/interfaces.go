package usecase

import (
	currency "p.solovev/internal/entity"
)

type (
	CurrencyReq interface {
		GetCurrencyRates(date string) (*currency.ValCurs, error)
		FindCurrencyRate(currency string, currencyRates *currency.ValCurs) (float64, error)
	}

	CurrencyContract interface {
		GetCurrencyRate(currency string, date string) (map[string]float64, error)
	}

	InfoContract interface {
		GetInfoAboutService() currency.InfoResponse
	}
)
