package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	apiGetAccountInformation = "/account"
	apiGetPositions          = "/positions"
	apiPostLeverage          = "/account/leverage"
	apiGetWalletBalances = "/wallet/balances"
)

type Account struct {
	client *Client
}

type AccountInformation struct {
	BackstopProvider             bool            `json:"backstopProvider"`
	Collateral                   decimal.Decimal `json:"collateral"`
	FreeCollateral               decimal.Decimal `json:"freeCollateral"`
	InitialMarginRequirement     decimal.Decimal `json:"initialMarginRequirement"`
	Liquidating                  bool            `json:"liquidating"`
	MaintenanceMarginRequirement decimal.Decimal `json:"maintenanceMarginRequirement"`
	MakerFee                     decimal.Decimal `json:"makerFee"`
	MarginFraction               decimal.Decimal `json:"marginFraction"`
	OpenMarginFraction           decimal.Decimal `json:"openMarginFraction"`
	TakerFee                     decimal.Decimal `json:"takerFee"`
	TotalAccountValue            decimal.Decimal `json:"totalAccountValue"`
	TotalPositionSize            decimal.Decimal `json:"totalPositionSize"`
	Username                     string          `json:"username"`
	Leverage                     decimal.Decimal `json:"leverage"`
	Positions                    []Position      `json:"positions"`
}

type Position struct {
	Cost                         decimal.Decimal `json:"cost"`
	EntryPrice                   decimal.Decimal `json:"entryPrice"`
	EstimatedLiquidationPrice    decimal.Decimal `json:"estimatedLiquidationPrice"`
	Future                       string          `json:"future"`
	InitialMarginRequirement     decimal.Decimal `json:"initialMarginRequirement"`
	LongOrderSize                decimal.Decimal `json:"longOrderSize"`
	MaintenanceMarginRequirement decimal.Decimal `json:"maintenanceMarginRequirement"`
	NetSize                      decimal.Decimal `json:"netSize"`
	OpenSize                     decimal.Decimal `json:"openSize"`
	RealizedPnl                  decimal.Decimal `json:"realizedPnl"`
	ShortOrderSize               decimal.Decimal `json:"shortOrderSize"`
	Side                         string          `json:"side"`
	Size                         decimal.Decimal `json:"size"`
	UnrealizedPnl                decimal.Decimal `json:"unrealizedPnl"`
	CollateralUsed               decimal.Decimal `json:"collateralUsed"`
}

func (a *Account) GetAccountInformation() (*AccountInformation, error) {
	request, err := a.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetAccountInformation),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := a.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *AccountInformation
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (a *Account) GetPositions() ([]Position, error) {
	request, err := a.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetPositions),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := a.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []Position
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (a *Account) ChangeAccountLeverage(leverage decimal.Decimal) error {
	body, err := json.Marshal(struct {
		Leverage decimal.Decimal `json:"leverage"`
	}{Leverage: leverage})
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := a.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiPostLeverage),
		Body:   body,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = a.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *Account) GetWalletBalances() ([]Balance, error) {
	request, err := a.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetWalletBalances),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := a.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []Balance
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}