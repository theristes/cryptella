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
	logger, _ := NewLogger()
	logger.Logf("Bought -> : %f at Price %f", c.amount, c.buyPrice)
}
