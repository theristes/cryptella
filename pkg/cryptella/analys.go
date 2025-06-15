package cryptella

func (c *Cryptella) analyze(current float64) {

	if !c.trade.canTrade() {
		c.Stop()
		return
	}

	c.trade.current = current
	if c.status == NONE {
		if c.trade.canBuy() {
			c.status = c.trade.placeOrder(c.api, BUY)
		}
	} else if c.status == BOUGHT {
		if c.trade.canSell() {
			c.status = c.trade.placeOrder(c.api, SELL)
		}
	} else if c.status == SOLD {
		c.trade.refresh(c.api)
		c.status = NONE
	}
	c.trade.logInfo(c.status)

}
