package cryptella

import (
	"log"
)

func (c *Cryptella) sell() {

	err := c.api.PlaceSellOrderOnApi(c.symbol, c.amount)
	if err != nil {
		log.Printf("Error placing sell order: %v", err)
		return
	}

	c.status = SOLD
	c.maxTrades--
	c.ShowInfo()

	log.Printf("Sold -> : %f at Price %f", c.amount, c.sellPrice)
	logger, _ := NewLogger()
	logger.Logf("Sold -> : %f at Price %f", c.amount, c.sellPrice)

}
