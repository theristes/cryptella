package cryptella

import (
	"github.com/samber/lo"
)

func (c *Cryptella) ShowMarketInfo() {

	values, _ := c.api.GetCandlesFromApi(c.symbol, c.interval, c.limit)
	println("Candles: ", len(values))

	find := values[len(values)-1].Close * (1 - (c.target + c.fee))

	candles := lo.Filter(values, func(item Candle, i int) bool {
		if item.Low <= find {
			return true
		}
		if item.Open <= find {
			return true
		}
		return false
	})

	println("Candles: ", len(candles))

}
