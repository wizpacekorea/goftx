package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID            int64           `json:"id"`
	Market        string          `json:"market"`
	Type          string          `json:"type"`
	Side          string          `json:"side"`
	Price         decimal.Decimal `json:"price"`
	Size          decimal.Decimal `json:"size"`
	FilledSize    decimal.Decimal `json:"filledSize"`
	RemainingSize decimal.Decimal `json:"remainingSize"`
	AvgFillPrice  decimal.Decimal `json:"avgFillPrice"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"createdAt"`
	ReduceOnly    bool            `json:"reduceOnly"`
	Ioc           bool            `json:"ioc"`
	PostOnly      bool            `json:"postOnly"`
	Future        string          `json:"future"`
	ClientID      string          `json:"clientId"`
}

type GetOrdersHistoryParams struct {
	Market    *string `json:"market"`
	Limit     *int    `json:"limit"`
	StartTime *int    `json:"start_time"`
	EndTime   *int    `json:"end_time"`
}

type TriggerOrder struct {
	ID               int64            `json:"id"`
	OrderID          int64            `json:"orderId"`
	Market           string           `json:"market"`
	CreatedAt        time.Time        `json:"createdAt"`
	Error            string           `json:"error"`
	Future           string           `json:"future"`
	OrderPrice       decimal.Decimal  `json:"orderPrice"`
	ReduceOnly       bool             `json:"reduceOnly"`
	Side             string           `json:"side"`
	Size             decimal.Decimal  `json:"size"`
	Status           string           `json:"status"`
	TrailStart       decimal.Decimal  `json:"trailStart"`
	TrailValue       decimal.Decimal  `json:"trailValue"`
	TriggerPrice     decimal.Decimal  `json:"triggerPrice"`
	TriggeredAt      time.Time        `json:"triggeredAt"`
	Type             string `json:"type"`
	OrderType        string        `json:"orderType"`
	FilledSize       decimal.Decimal  `json:"filledSize"`
	AvgFillPrice     decimal.Decimal  `json:"avgFillPrice"`
	OrderStatus      string           `json:"orderStatus"`
	RetryUntilFilled bool             `json:"retryUntilFilled"`
}

type GetOpenTriggerOrdersParams struct {
	Market *string           `json:"market"`
	Type   *string `json:"type"`
}

type Trigger struct {
	Error      string    `json:"error"`
	FilledSize float64   `json:"filledSize"`
	OrderSize  float64   `json:"orderSize"`
	OrderID    int64     `json:"orderId"`
	Time       time.Time `json:"time"`
}

type GetTriggerOrdersHistoryParams struct {
	Market    *string           `json:"market"`
	StartTime *int              `json:"start_time"`
	EndTime   *int              `json:"end_time"`
	Side      *string           `json:"side"`
	Type      *string `json:"type"`
	OrderType *string        `json:"orderType"`
	Limit     *int              `json:"limit"`
}

type PlaceOrderPayload struct {
	Market                  string          `json:"market"`
	Side                    string          `json:"side"`
	Price                   decimal.Decimal `json:"price"`
	Type                    string          `json:"type"`
	Size                    decimal.Decimal `json:"size"`
	ReduceOnly              bool            `json:"reduceOnly,omitempty"`
	IOC                     bool            `json:"ioc,omitempty"`
	PostOnly                bool            `json:"postOnly,omitempty"`
	ClientID                string          `json:"clientId,omitempty"`
	ExternalReferralProgram string          `json:"externalReferralProgram,omitempty"`
}

type PlaceTriggerOrderPayload struct {
	Market           string           `json:"market"`
	Side             string           `json:"side"`
	Size             decimal.Decimal  `json:"size"`
	Type             string `json:"type"`
	ReduceOnly       bool             `json:"reduceOnly,omitempty"`
	RetryUntilFilled bool             `json:"retryUntilFilled,omitempty"`
	TriggerPrice     *decimal.Decimal `json:"triggerPrice,omitempty"`
	OrderPrice       *decimal.Decimal `json:"orderPrice,omitempty"`
	TrailValue       *decimal.Decimal `json:"trailValue,omitempty"`
}

func (t PlaceTriggerOrderPayload) Validate() error {
	switch t.Type {
	case TriggerTypeStop:
		if t.TriggerPrice == nil {
			return errors.New("triggerPrice is required for stop loss orders")
		}
	case TriggerTypeTrailingStop:
		if t.TrailValue == nil {
			return errors.New("trailValue is required for trailing stop orders")
		}
	case TriggerTypeTakeProfit:
		if t.TriggerPrice == nil {
			return errors.New("triggerPrice is required for take profit orders")
		}
	}

	return nil
}

type ModifyOrderPayload struct {
	Price    *decimal.Decimal `json:"price,omitempty"`
	Size     *decimal.Decimal `json:"size,omitempty"`
	ClientID *string          `json:"clientId,omitempty"`
}

type ModifyTriggerOrderPayload struct {
	Size         decimal.Decimal  `json:"size"`
	TriggerPrice decimal.Decimal  `json:"triggerPrice"`
	OrderPrice   *decimal.Decimal `json:"orderPrice,omitempty"`
	TrailValue   *decimal.Decimal `json:"trailValue,omitempty"`
}

type CancelAllOrdersPayload struct {
	Market                *string `json:"market,omitempty"`
	ConditionalOrdersOnly *bool   `json:"conditionalOrdersOnly,omitempty"`
	LimitOrdersOnly       *bool   `json:"limitOrdersOnly"`
}

const (
	apiOrders                  = "/orders"
	apiGetOrdersHistory        = "/orders/history"
	apiModifyOrder             = "/orders/%d/modify"
	apiModifyOrderByClientID   = "/orders/by_client_id/%d/modify"
	apiTriggerOrders           = "/conditional_orders"
	apiGetOrderTriggers        = "/conditional_orders/%d/triggers"
	apiGetTriggerOrdersHistory = "/conditional_orders/history"
	apiModifyTriggerOrder      = "/conditional_orders/%d/modify"
)

type Orders struct {
	client *Client
}

func (o *Orders) GetOpenOrders(market string) ([]Order, error) {
	requestParams := Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiOrders),
	}
	if market != "" {
		requestParams.Params = map[string]string{
			"market": market,
		}
	}

	request, err := o.client.prepareRequest(requestParams)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOrdersHistory(params *GetOrdersHistoryParams) ([]Order, error) {
	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetOrdersHistory),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOpenTriggerOrders(params *GetOpenTriggerOrdersParams) ([]TriggerOrder, error) {
	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiTriggerOrders),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []TriggerOrder
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOrderTriggers(orderID int64) ([]Trigger, error) {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, fmt.Sprintf(apiGetOrderTriggers, orderID)),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []Trigger
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetTriggerOrdersHistory(params *GetTriggerOrdersHistoryParams) ([]TriggerOrder, error) {
	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetTriggerOrdersHistory),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []TriggerOrder
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) PlaceOrder(payload *PlaceOrderPayload) (*Order, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiOrders),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) PlaceTriggerOrder(payload *PlaceTriggerOrderPayload) (*TriggerOrder, error) {
	err := payload.Validate()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiTriggerOrders),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *TriggerOrder
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) ModifyOrder(payload *ModifyOrderPayload, orderID int64) (*Order, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, fmt.Sprintf(apiModifyOrder, orderID)),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) ModifyOrderByClientID(payload *ModifyOrderPayload, clientOrderID int64) (*Order, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, fmt.Sprintf(apiModifyOrderByClientID, clientOrderID)),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) ModifyTriggerOrder(payload *ModifyTriggerOrderPayload, orderID int64) (*TriggerOrder, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, fmt.Sprintf(apiModifyTriggerOrder, orderID)),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *TriggerOrder
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOrder(orderID int64) (*Order, error) {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s/%d", apiUrl, apiOrders, orderID),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOrderByClientID(clientOrderID int64) (*Order, error) {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s/by_client_id/%d", apiUrl, apiOrders, clientOrderID),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) CancelOrder(orderID int64) error {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s/%d", apiUrl, apiOrders, orderID),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (o *Orders) CancelOrderByClientID(clientOrderID int64) error {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s/by_client_id/%d", apiUrl, apiOrders, clientOrderID),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (o *Orders) CancelOpenTriggerOrder(triggerOrderID int64) error {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s/%d", apiUrl, apiTriggerOrders, triggerOrderID),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (o *Orders) CancelAllOrders(payload *CancelAllOrdersPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiOrders),
		Body:   body,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
