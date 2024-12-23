package cryptella

import (
	"github.com/samber/lo"
)

func (c Cryptella) ShowMarketInfo() {

	interval := "1m"
	limit := 60
	values, _ := c.api.GetCandlesFromApi(c.trade.symbol, interval, limit)
	println("Full Candles: ", len(values))

	find := basedUpOnLastValue(values, c)

	candles := c.filterCandles(values, find)

	println("Last Value Candles: ", len(candles))

	find = basedUpOnAverageValue(values, c)

	candles = c.filterCandles(values, find)

	println("Average Value Candles: ", len(candles))

}

func (c Cryptella) filterCandles(values []Candle, find float64) []Candle {

	candles := lo.Filter(values, func(item Candle, i int) bool {
		if item.Low <= find {
			return true
		}
		if item.Open <= find {
			return true
		}
		return false
	})

	return candles
}

func basedUpOnLastValue(values []Candle, c Cryptella) float64 {
	return values[len(values)-1].Close * (1 - (c.trade.target + c.trade.fee))
}

func basedUpOnAverageValue(values []Candle, c Cryptella) float64 {
	var sum float64
	for _, value := range values {
		sum += value.Close
	}

	discount := 1 - (c.trade.target + c.trade.fee)
	media := sum / float64(len(values))
	return media * discount

}
