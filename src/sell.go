package cryptella

import (
	"log"
)

func (c *Cryptella) sell() {

	err := c.api.PlaceSellOrderOnApi(c.symbol, c.amount, c.sellPrice)
	if err != nil {
		log.Printf("Error placing sell order: %v", err)
		return
	}

	c.maxTrades--
	c.ShowInfo()

	logger, _ := NewLogger()
	logger.Logf("Sold -> : %f at Price %f", c.amount, c.sellPrice)

}
