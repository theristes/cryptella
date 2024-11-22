package cryptella

import (
	"log"
)

func (c *Cryptella) buy() {

	err := c.api.PlaceBuyOrderOnApi(c.symbol, c.amount)
	if err != nil {
		log.Printf("Error placing buy order: %v", err)
		return
	}

	c.amount = c.amount * (1 - c.fee)
	c.status = BOUGHT

	log.Printf("Bought -> : %f at Price %f", c.amount, c.buyPrice)
	logger, _ := NewLogger()
	logger.Logf("Bought -> : %f at Price %f", c.amount, c.buyPrice)
}
