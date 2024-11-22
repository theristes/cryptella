package cryptella

import "math"

func (c *Cryptella) fillAmount() {

	value, _ := c.api.GetBalanceFreeApi("USDT")

	available := value * 0.75

	price, _, _, _, _ := c.api.GetPriceFromApi(c.symbol)

	c.amount = c.calculateAmount(available, price)
}

func (c *Cryptella) calculateAmount(available float64, price float64) float64 {

	qty := available / price

	minQty, maxQty, stepSize, minNotional, _, _, _, _, _, err := c.api.GetFiltersFromApi(c.symbol)
	if err != nil {
		return 0
	}

	if qty < minQty {
		qty = minQty
	}
	if qty > maxQty {
		qty = maxQty
	}

	qty = math.Floor(qty/stepSize) * stepSize

	if qty*stepSize*price < minNotional {

		qty = math.Ceil(minNotional/price/stepSize) * stepSize

	}

	return qty
}
