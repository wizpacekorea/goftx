package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type SubAccount struct {
	Nickname    string `json:"nickname"`
	Deletable   bool   `json:"deletable"`
	Editable    bool   `json:"editable"`
	Competition bool   `json:"competition,omitempty"`
}

type Balance struct {
	Coin  string          `json:"coin"`
	Free  decimal.Decimal `json:"free"`
	Total decimal.Decimal `json:"total"`
}

type TransferPayload struct {
	Coin        string          `json:"coin"`
	Size        decimal.Decimal `json:"size"`
	Source      *string         `json:"source"`
	Destination *string         `json:"destination"`
}

type TransferResponse struct {
	ID     int64           `json:"id"`
	Coin   string          `json:"coin"`
	Size   decimal.Decimal `json:"size"`
	Time   time.Time       `json:"time"`
	Notes  string          `json:"notes"`
	Status string  `json:"status"`
}

const (
	apiSubAccounts           = "/subaccounts"
	apiChangeSubAccountName  = "/subaccounts/update_name"
	apiGetSubAccountBalances = "/subaccounts/%s/balances"
	apiTransfer              = "/subaccounts/transfer"
)

type SubAccounts struct {
	client *Client
}

func (s *SubAccounts) GetSubAccounts() ([]SubAccount, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiSubAccounts),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []SubAccount
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SubAccounts) CreateSubaccount(nickname string) (*SubAccount, error) {
	body, err := json.Marshal(struct {
		Nickname string `json:"nickname"`
	}{Nickname: nickname})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiSubAccounts),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result SubAccount
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}

func (s *SubAccounts) ChangeSubaccount(nickname, newNickname string) error {
	body, err := json.Marshal(struct {
		Nickname    string `json:"nickname"`
		NewNickname string `json:"newNickname"`
	}{Nickname: nickname, NewNickname: newNickname})
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiChangeSubAccountName),
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

func (s *SubAccounts) DeleteSubAccount(nickname string) error {
	body, err := json.Marshal(struct {
		Nickname string `json:"nickname"`
	}{Nickname: nickname})
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiSubAccounts),
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

func (s *SubAccounts) GetSubAccountBalances(nickname string) ([]Balance, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, fmt.Sprintf(apiGetSubAccountBalances, nickname)),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
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

func (s *SubAccounts) Transfer(payload *TransferPayload) (*TransferResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiTransfer),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result TransferResponse
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}
