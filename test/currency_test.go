package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	currency "p.solovev/internal/entity"
	"p.solovev/internal/usecase"
	"strconv"
	"strings"
	"testing"
)

type mockCurrencyReq struct{}

func (m *mockCurrencyReq) GetCurrencyRates(date string) (*currency.ValCurs, error) {
	ValCurs := new(currency.ValCurs)
	valute := new(currency.Valute)
	valute.CharCode = "USD"
	valute.Value = "30,15"
	ValCurs.Valutes = append(ValCurs.Valutes, *valute)
	return ValCurs, nil
}

func (m *mockCurrencyReq) FindCurrencyRate(currency string, currencyRates *currency.ValCurs) (float64, error) {
	for _, v := range currencyRates.Valutes {
		if v.CharCode == currency {
			return strconv.ParseFloat(strings.Replace(v.Value, ",", ".", 1), 64)
		}
	}
	return 0, fmt.Errorf("currency or rate not found")
}

func TestGetCurrencyRate(t *testing.T) {
	mockCBRF := &mockCurrencyReq{}
	useCase := usecase.NewCurrencyUseCase(mockCBRF)

	currencyCode := "USD"
	date := ""

	result, err := useCase.GetCurrencyRate(currencyCode, date)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
