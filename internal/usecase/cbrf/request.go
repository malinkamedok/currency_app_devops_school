package cbrf

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
	logger2 "p.solovev/pkg/logger"

	"golang.org/x/net/html/charset"
	currency "p.solovev/internal/entity"
	"p.solovev/internal/usecase"
)

type CurrencyCBRF struct{}

func NewCurrencyReq() *CurrencyCBRF {
	return &CurrencyCBRF{}
}

var _ usecase.CurrencyReq = (*CurrencyCBRF)(nil)

func (i CurrencyCBRF) initRequest(dateFormatted string) (*http.Request, error) {
	url := "https://www.cbr.ru/scripts/XML_daily.asp?date_req=" + dateFormatted

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger2.Error("Error in creating request", zap.Error(err))
		return nil, err
	}

	// Getting 403 Forbidden error without setting this header
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")

	return req, nil
}

func (i CurrencyCBRF) sendRequest(r *http.Request) (*http.Response, error) {
	c := http.Client{}

	resp, err := c.Do(r)
	if err != nil {
		logger2.Error("Error in sending request", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (i CurrencyCBRF) decodeResponse(response *http.Response) (*currency.ValCurs, error) {
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logger2.Error("Error in reading response", zap.Error(err))
		return nil, err
	}

	reader := bytes.NewReader(responseData)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	rates := new(currency.ValCurs)

	err = decoder.Decode(rates)
	if err != nil {
		logger2.Error("Error in decoding response", zap.Error(err))
		return nil, err
	}

	return rates, nil
}

func (i CurrencyCBRF) GetCurrencyRates(date string) (*currency.ValCurs, error) {
	req, err := i.initRequest(date)
	if err != nil {
		return nil, err
	}

	resp, err := i.sendRequest(req)
	if err != nil {
		return nil, err
	}

	currencyRates, err := i.decodeResponse(resp)
	if err != nil {
		return nil, err
	}

	return currencyRates, nil
}

func (i CurrencyCBRF) FindCurrencyRate(currency string, currencyRates *currency.ValCurs) (float64, error) {
	for _, v := range currencyRates.Valutes {
		if v.CharCode == currency {
			return strconv.ParseFloat(strings.Replace(v.Value, ",", ".", 1), 64)
		}
	}
	return 0, fmt.Errorf("Currency or rate not found")
}
