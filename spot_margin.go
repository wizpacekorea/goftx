package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type BorrowRate struct {
	Coin     string  `json:"coin"`
	Estimate float64 `json:"estimate"`
	Previous float64 `json:"previous"`
}

type LendingRate BorrowRate

type BorrowSummary struct {
	Coin string          `json:"coin"`
	Size decimal.Decimal `json:"size"`
}

type SpotMarginMarketInfo struct {
	Coin         string          `json:"coin"`
	Borrowed     decimal.Decimal `json:"borrowed"`
	Free         decimal.Decimal `json:"free"`
	EstimatedRate float64         `json:"estimatedRate"`
	PreviousRate float64         `json:"previousRate"`
}

type GetSpotMarginMarketInfoResponse struct {
	Base  SpotMarginMarketInfo
	Quote SpotMarginMarketInfo
}

type BorrowHistory struct {
	Coin string          `json:"coin"`
	Cost decimal.Decimal `json:"cost"`
	Rate decimal.Decimal `json:"rate"`
	Size decimal.Decimal `json:"size"`
	Time time.Time       `json:"time"`
}

type LendingHistory BorrowHistory

type LendingOffer struct {
	Coin string          `json:"coin"`
	Rate float64         `json:"rate"`
	Size decimal.Decimal `json:"size"`
}

type LendingInfo struct {
	Coin     string          `json:"coin"`
	Lendable decimal.Decimal `json:"lendable"`
	Locked   decimal.Decimal `json:"locked"`
	MinRate  float64         `json:"minRate"`
	Offered  decimal.Decimal `json:"offered"`
}

type LendingOfferPayload struct {
	Coin string          `json:"coin"`
	Size decimal.Decimal `json:"size"`
	Rate float64         `json:"rate"`
}

const (
	apiBorrowRates    = "/spot_margin/borrow_rates"
	apiLendingRates   = "/spot_margin/lending_rates"
	apiBorrowSummary  = "/spot_margin/borrow_summary"
	apiMarketInfo     = "/spot_margin/market_info"
	apiBorrowHistory  = "/spot_margin/borrow_history"
	apiLendingHistory = "/spot_margin/lending_history"
	apiLendingOffers  = "/spot_margin/offers"
	apiLendingInfo    = "/spot_margin/lending_info"
)

type SpotMargin struct {
	client *Client
}

func (s *SpotMargin) GetBorrowRates() ([]BorrowRate, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiBorrowRates),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []BorrowRate
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingRates() ([]LendingRate, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiLendingRates),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []LendingRate
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetDailyBorrowedAmounts() ([]BorrowSummary, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiBorrowSummary),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []BorrowSummary
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetMarketInfo(market string) ([]GetSpotMarginMarketInfoResponse, error) {
	queryParams := map[string]string{
		"market": market,
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiMarketInfo),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []GetSpotMarginMarketInfoResponse
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetBorrowHistory() ([]BorrowHistory, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiBorrowHistory),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []BorrowHistory
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingHistory() ([]LendingHistory, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiLendingHistory),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []LendingHistory
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingOffers() ([]LendingOffer, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiLendingOffers),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []LendingOffer
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingInfo() ([]LendingInfo, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiLendingInfo),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []LendingInfo
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) SubmitLendingOffer(payload *LendingOfferPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiLendingOffers),
		Body:   body,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = s.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
