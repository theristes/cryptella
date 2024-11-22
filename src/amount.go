package cryptella

import (
	"math"
)

func (c *Cryptella) fillAmount() {

	value, _ := c.api.GetBalanceFreeApi("USDT")

	available := value * c.balanceAsset

	price, _, _, _, _ := c.api.GetPriceFromApi(c.symbol)

	qty := available / price
	c.amount = c.calculateAmount(qty)
}

func (c *Cryptella) calculateAmount(qty float64) float64 {

	minQty, maxQty, stepSize, _, _, _, _, _, _, err := c.api.GetFiltersFromApi(c.symbol)
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

	return qty
}
