package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const apiFills = "/fills"

type Fills struct {
	client *Client
}

type GetFillsParams struct {
	Market    *string `json:"market"`
	Limit     *int    `json:"limit"`
	StartTime *int64  `json:"start_time"`
	EndTime   *int64  `json:"end_time"`
	Order     *string `json:"order"`
	OrderID   *int64  `json:"orderId"`
}

type Fill struct {
	Fee           float64         `json:"fee"`
	FeeCurrency   string          `json:"feeCurrency"`
	FeeRate       float64         `json:"feeRate"`
	Future        string          `json:"future"`
	ID            int64           `json:"id"`
	Liquidity     string       `json:"liquidity"`
	Market        string          `json:"market"`
	BaseCurrency  string          `json:"baseCurrency"`
	QuoteCurrency string          `json:"quoteCurrency"`
	OrderID       int64           `json:"orderId"`
	TradeID       int64           `json:"tradeId"`
	Price         decimal.Decimal `json:"price"`
	Side          string            `json:"side"`
	Size          decimal.Decimal `json:"size"`
	Time          FTXTime         `json:"time"`
	Type          string          `json:"type"`
}

func (f *Fills) GetFills(params *GetFillsParams) ([]Fill, error) {
	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := f.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiFills),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := f.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []Fill
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
