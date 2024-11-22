package cryptella

import (
	"fmt"
)

func (c *Cryptella) analyze() {

	// Fetch the current market price
	price, _, _, _, _ := c.api.GetPriceFromApi(c.symbol)
	fmt.Printf("Price: %f, Need: %f\n", price, c.lowPrice)

	switch c.status {
	case "NONE":

		if c.buyStrategy(price) {
			c.status = "BUY"
			c.buyPrice = price
			fmt.Printf("BUY signal detected at price %.2f\n", price)
		}
	case "BOUGHT":

		if c.sellStrategy(price) {
			c.status = "SELL"
			c.sellPrice = price
			fmt.Printf("SELL signal detected at price %.2f\n", price)
		}
	default:
		fmt.Println("No action for current status:", c.status)
	}
}

func (c *Cryptella) sellStrategy(price float64) bool {
	return price >= c.buyPrice*(1+c.target+c.fee)
}

func (c *Cryptella) buyStrategy(price float64) bool {
	return price <= c.lowPrice
}

func (c *Cryptella) fillLowPrice() {

	discount := 1 - ((c.target + c.fee) * c.buyThresehold) // 15% of target and fee to find a good price

	values, _ := c.api.GetCandlesFromApi(c.symbol, c.interval, c.limit)

	// get the media from values
	var sum float64
	for _, value := range values {
		sum += value.Close
	}
	media := sum / float64(len(values))

	c.lowPrice = media * discount
}
