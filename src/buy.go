package cryptella

import (
	"log"
)

func (c *Cryptella) buy() {

	err := c.api.PlaceBuyOrderOnApi(c.symbol, c.amount, c.buyPrice)
	if err != nil {
		log.Printf("Error placing buy order: %v", err)
		return
	}
	c.status = BOUGHT

	c.fillSellPrice() // We need to fill the sell price after buying

	log.Printf("Bought -> : %f at Price %f", c.amount, c.buyPrice)
	logger, _ := NewLogger()
	logger.Logf("Bought -> : %f at Price %f", c.amount, c.buyPrice)
}
