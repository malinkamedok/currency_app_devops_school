package usecase

import (
	"fmt"
	"strings"
	"time"

	logger2 "p.solovev/pkg/logger"
)

type CurrencyUseCase struct {
	cbrf CurrencyReq
}

var _ CurrencyContract = (*CurrencyUseCase)(nil)

func NewCurrencyUseCase(cbrf CurrencyReq) *CurrencyUseCase {
	return &CurrencyUseCase{cbrf: cbrf}
}

func (i CurrencyUseCase) GetCurrencyRate(currencyCode string, date string) (map[string]float64, error) {

	dateFormatted, err := parseAndFormatDate(date)
	if err != nil {
		return nil, err
	}

	rates, err := i.cbrf.GetCurrencyRates(dateFormatted)
	if err != nil {
		return nil, err
	}

	currencyResult := make(map[string]float64)

	if currencyCode == "" {
		for _, v := range rates.Valutes {
			currencyRate, err := i.cbrf.FindCurrencyRate(v.CharCode, rates)
			if err != nil {
				return nil, err
			}

			currencyResult[v.CharCode] = currencyRate
		}
	} else {
		currencyCode = strings.ToUpper(currencyCode)

		currencyRate, err := i.cbrf.FindCurrencyRate(currencyCode, rates)
		if err != nil {
			return nil, err
		}

		currencyResult[currencyCode] = currencyRate
	}

	return currencyResult, nil
}

func parseAndFormatDate(date string) (string, error) {
	var dateFormatted string

	if date == "" {
		dateFormatted = time.Now().Format("02/01/2006")
	} else {
		dateParsed, err := time.Parse("2006-01-02", date)
		if err != nil {
			logger2.Error("Error in parsing date")
			return "", err
		}

		if dateParsed.After(time.Now()) {
			return "", fmt.Errorf("incorrect time")
		}

		dateFormatted = dateParsed.Format("02/01/2006")
	}

	return dateFormatted, nil
}
