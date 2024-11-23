package cryptella

import (
	"fmt"
	"time"
)

func (c *Cryptella) wait() {
	time.Sleep(500 * time.Millisecond)
}
func (c *Cryptella) analyze() {

	// Fetch the current market price
	price, _, _, _, _ := c.api.GetPriceFromApi(c.symbol)

	if c.status == NONE {
		fmt.Printf("Price: %f, Need: %f\n", price, c.getLowPrice())
		c.wait()
		if c.buyStrategy(price) {
			c.status = BUY
			c.buyPrice = price
			fmt.Printf("BUY signal detected at price %f\n", price)
		}
	} else if c.status == BUY {
		c.status = BOUGHT
		c.buy()
		c.fillSellPrices()
		fmt.Printf("Bought -> : %f at Price %f", c.amount, c.buyPrice)

	} else if c.status == BOUGHT {
		fmt.Printf("Price: %f, Bought: %f, Need Sell: %f\n", price, c.buyPrice, c.sellPrice)
		c.wait()
		if c.sellStrategy(price) {
			c.status = SELL
			c.sellPrice = price
			fmt.Printf("SELL signal detected at price %f\n", price)
		}

	} else if c.status == SELL {
		c.status = SOLD
		c.sell()
		fmt.Printf("Sold -> : %f at Price %f", c.amount, c.sellPrice)
	} else if c.status == SOLD {
		c.status = NONE
		c.fillLowPrice()
		fmt.Printf("Trade completed\n")
	}

}

func (c *Cryptella) sellStrategy(price float64) bool {

	if price <= c.buyPrice*(1-c.stopLoss) {

		logger, _ := NewLogger()
		logger.Logf("Stop loss triggered at price %f", price)
		fmt.Printf("Stop loss triggered at price %f\n", price)

		return true
	}

	return price >= c.sellPrice
}

func (c *Cryptella) buyStrategy(price float64) bool {
	return price <= c.getLowPrice()
}

func (c *Cryptella) getLowPrice() float64 {
	return c.lowPrice * (1 - c.getBuyThresehold())
}

func (c *Cryptella) getBuyThresehold() float64 {
	c.buyThresehold = c.buyThresehold - 0.01
	return c.buyThresehold
}

func (c *Cryptella) fillLowPrice() {
	discount := 1 - (c.target + c.fee) // 15% of target and fee to find a good price

	values, _ := c.api.GetCandlesFromApi(c.symbol, c.interval, c.limit)

	// get the media from values
	var sum float64
	for _, value := range values {
		sum += value.Close
	}
	media := sum / float64(len(values))

	c.lowPrice = media * discount
}

func (c *Cryptella) fillSellPrices() {
	c.sellPrice = c.buyPrice * (1 + c.target + c.fee)
}
