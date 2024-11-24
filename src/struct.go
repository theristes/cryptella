package cryptella

import (
	"time"

	"github.com/adshao/go-binance/v2"
)

type Cryptella struct {
	api    *Api
	trade  *Trade
	status string // BUY |  BOUGHT | SELL | SOLD | NONE
	stopCh chan struct{}
}

type Trade struct {
	symbol          string
	amount          float64
	fee             float64
	target          float64
	stopLoss        float64
	reach           float64
	current         float64
	reached         float64
	limitLossOrders int64
}

type Api struct {
	client *binance.Client
}

type Candle struct {
	OpenTime  time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	CloseTime time.Time
	QAV       float64
	TradeNum  int64
	TBBV      float64
	TBQAV     float64
}

const (
	BUY    = "BUY"
	BOUGHT = "BOUGHT"
	SELL   = "SELL"
	SOLD   = "SOLD"
	NONE   = "NONE"
)
