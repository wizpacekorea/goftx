package goftx

import (
	"encoding/json"
	"math"
	"time"
)

type Resolution int

const (
	Sec15    = 15
	Minute   = 60
	Minute5  = 300
	Minute15 = 900
	Hour     = 3600
	Hour4    = 14400
	Day      = 86400
)

const (
	OrderBookChannel = "orderbook"
	TradesChannel    = "trades"
	TickerChannel    = "ticker"
	MarketsChannel   = "markets"
	FillsChannel     = "fills"
	OrdersChannel    = "orders"
)

const (
	Subscribe   = "subscribe"
	UnSubscribe = "unsubscribe"
	Login       = "login"
)

const (
	ResponseTypeError        = "error"
	ResponseTypeSubscribed   = "subscribed"
	ResponseTypeUnSubscribed = "unsubscribed"
	ResponseTypeInfo         = "info"
	ResponseTypePartial      = "partial"
	ResponseTypeUpdate       = "update"
)

const TransferStatusComplete = "complete"
const (
	OrderTypeLimitOrder  = "limit"
	OrderTypeMarketOrder = "market"
)

const (
	SideSell = "sell"
	SideBuy  = "buy"
)

const (
	OrderStatusNew    = "new"
	OrderStatusOpen   = "open"
	OrderStatusClosed = "closed"
)

const (
	TriggerTypeStop         = "stop"
	TriggerTypeTrailingStop = "trailingStop"
	TriggerTypeTakeProfit   = "takeProfit"
)

type FTXTime struct {
	Time time.Time
}

func (f *FTXTime) UnmarshalJSON(data []byte) error {
	var t float64
	err := json.Unmarshal(data, &t)

	// FTX uses ISO format sometimes so we have to detect and handle that differently.
	if err != nil {
		var iso time.Time
		errIso := json.Unmarshal(data, &iso)

		if errIso != nil {
			return err
		}

		f.Time = iso
		return nil
	}

	sec, nsec := math.Modf(t)
	f.Time = time.Unix(int64(sec), int64(nsec))
	return nil
}

func (f FTXTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f.Time.UnixNano()) / float64(1000000000))
}

const (
	LiquidityTaker = "taker"
	LiquidityMaker = "maker"
)

const (
	FutureTypeFuture    = "future"
	FutureTypePerpetual = "perpetual"
	FutureTypeMove      = "move"
)
