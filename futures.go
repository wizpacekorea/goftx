package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Future struct {
	Ask                 decimal.Decimal `json:"ask"`
	Bid                 decimal.Decimal `json:"bid"`
	Change1h            decimal.Decimal `json:"change1h"`
	Change24h           decimal.Decimal `json:"change24h"`
	ChangeBod           decimal.Decimal `json:"changeBod"`
	VolumeUsd24h        float64         `json:"volumeUsd24h"`
	Volume              float64         `json:"volume"`
	Description         string          `json:"description"`
	Enabled             bool            `json:"enabled"`
	Expired             bool            `json:"expired"`
	Expiry              time.Time       `json:"expiry"`
	Index               float64         `json:"index"`
	ImfFactor           float64         `json:"imfFactor"`
	Last                decimal.Decimal `json:"last"`
	LowerBound          decimal.Decimal `json:"lowerBound"`
	Mark                decimal.Decimal `json:"mark"`
	Name                string          `json:"name"`
	Perpetual           bool            `json:"perpetual"`
	PositionLimitWeight float64         `json:"positionLimitWeight"`
	PostOnly            bool            `json:"postOnly"`
	PriceIncrement      decimal.Decimal `json:"priceIncrement"`
	SizeIncrement       decimal.Decimal `json:"sizeIncrement"`
	Underlying          string          `json:"underlying"`
	UpperBound          decimal.Decimal `json:"upperBound"`
	Type                string      `json:"type"`
}

type FutureExpired struct {
	Ask                   decimal.Decimal `json:"ask"`
	Bid                   decimal.Decimal `json:"bid"`
	Description           string          `json:"description"`
	Enabled               bool            `json:"enabled"`
	Expired               bool            `json:"expired"`
	Expiry                time.Time       `json:"expiry"`
	ExpiryDescription     string          `json:"expiryDescription"`
	Group                 string          `json:"group"`
	ImfFactor             decimal.Decimal `json:"imfFactor"`
	Index                 decimal.Decimal `json:"index"`
	Last                  decimal.Decimal `json:"last"`
	LowerBound            decimal.Decimal `json:"lowerBound"`
	MarginPrice           decimal.Decimal `json:"marginPrice"`
	Mark                  decimal.Decimal `json:"mark"`
	MoveStart             string          `json:"moveStart"`
	Name                  string          `json:"name"`
	Perpetual             bool            `json:"perpetual"`
	PositionLimitWeight   decimal.Decimal `json:"positionLimitWeight"`
	PostOnly              bool            `json:"postOnly"`
	PriceIncrement        decimal.Decimal `json:"priceIncrement"`
	SizeIncrement         decimal.Decimal `json:"sizeIncrement"`
	Type                  string          `json:"type"`
	Underlying            string          `json:"underlying"`
	UnderlyingDescription string          `json:"underlyingDescription"`
	UpperBound            decimal.Decimal `json:"upperBound"`
}

type FutureStats struct {
	Volume                   decimal.Decimal `json:"volume"`
	NextFundingRate          float64         `json:"nextFundingRate"`
	NextFundingTime          time.Time       `json:"nextFundingTime"`
	ExpirationPrice          decimal.Decimal `json:"expirationPrice"`
	PredictedExpirationPrice decimal.Decimal `json:"predictedExpirationPrice"`
	StrikePrice              decimal.Decimal `json:"strikePrice"`
	OpenInterest             float64         `json:"openInterest"`
}

type GetFundingRatesParams struct {
	StartTime *int    `json:"start_time"`
	EndTime   *int    `json:"end_time"`
	Future    *string `json:"future"`
}

type FundingRate struct {
	Future string          `json:"future"`
	Rate   decimal.Decimal `json:"rate"`
	Time   time.Time       `json:"time"`
}

type GetHistoricalIndexParams struct {
	IndexName  string `json:"index_name"`
	Resolution int    `json:"resolution"`
	Limit      *int   `json:"limit"`
	StartTime  *int   `json:"start_time"`
	EndTime    *int   `json:"end_time"`
}

type HistoricalIndex struct {
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	StartTime time.Time       `json:"startTime"`
	Volume    decimal.Decimal `json:"volume"`
}

const (
	apiFutures        = "/futures"
	apiFundingRates   = "/funding_rates"
	apiIndexWeights   = "/indexes/%s/weights"
	apiIndexCandles   = "/indexes/%s/candles"
	apiExpiredFutures = "expired_futures"
)

type Futures struct {
	client *Client
}

func (f *Futures) GetFutures() ([]Future, error) {
	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiFutures),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []Future
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (f *Futures) GetFuture(name string) (*Future, error) {
	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s/%s", apiUrl, apiFutures, name),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *Future
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (f *Futures) GetFutureStats(name string) (*FutureStats, error) {
	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s/%s/stats", apiUrl, apiFutures, name),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *FutureStats
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (f *Futures) GetFundingRates(params *GetFundingRatesParams) ([]FundingRate, error) {
	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiFundingRates),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []FundingRate
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (f *Futures) GetIndexWeights(indexName string) (map[string]decimal.Decimal, error) {
	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, fmt.Sprintf(apiIndexWeights, indexName)),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make(map[string]decimal.Decimal)
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (f *Futures) GetExpiredFutures() ([]FutureExpired, error) {
	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiExpiredFutures),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []FutureExpired
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (f *Futures) GetHistoricalIndex(market string, params *GetHistoricalIndexParams) ([]HistoricalIndex, error) {
	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := f.client.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, fmt.Sprintf(apiIndexCandles, market)),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []HistoricalIndex
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
