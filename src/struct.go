package cryptella

import (
	"time"

	"github.com/adshao/go-binance/v2"
)

type Cryptella struct {
	api Api

	symbol    string
	amount    float64
	buyPrice  float64
	lowPrice  float64
	sellPrice float64

	// Config
	fee           float64
	target        float64
	maxTrades     int
	interval      string
	limit         int
	buyThresehold float64

	status string // BUY |  BOUGHT | SELL | SOLD | NONE
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
