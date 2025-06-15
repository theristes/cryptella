package cryptella

import (
	"fmt"
	"os"
	"strconv"
)

func NewTrade(api *Api) *Trade {

	symbol := os.Getenv("SYMBOL")
	amount, _ := strconv.ParseFloat(os.Getenv("AMOUNT"), 64)
	fee, _ := strconv.ParseFloat(os.Getenv("FEE"), 64)
	target, _ := strconv.ParseFloat(os.Getenv("TARGET"), 64)
	stopLoss, _ := strconv.ParseFloat(os.Getenv("STOP_LOSS"), 64)
	limitLossOrders, _ := strconv.ParseInt(os.Getenv("LIMIT_LOSS_ORDERS"), 10, 64)
	simplified, _ := strconv.ParseBool(os.Getenv("SIMPLIFIED"))

	discount := 1 - (target + fee)
	media, _ := api.GetMediaFromApi(symbol)

	var reach float64
	if simplified {
		reach = getSimplifiedBuyPrice()
		if reach == -1 {
			reach = media * discount
		}
	} else {
		reach = media * discount
	}
	current, _ := api.GetPriceFromApi(symbol)

	return &Trade{
		symbol:          symbol,
		amount:          amount,
		fee:             fee,
		target:          target,
		stopLoss:        stopLoss,
		reach:           reach,
		current:         current,
		limitLossOrders: limitLossOrders,
		simplified:      simplified,
	}
}

func getSimplifiedBuyPrice() float64 {
	buyPrice, _ := strconv.ParseFloat(os.Getenv("BUY_PRICE"), 64)
	return buyPrice
}

func getSimplifiedSellPrice() float64 {
	buyPrice, _ := strconv.ParseFloat(os.Getenv("SELL_PRICE"), 64)
	return buyPrice
}

func (t Trade) canTrade() bool {
	return t.limitLossOrders > 0
}

func (t Trade) logInfo(status string) {
	logger, _ := NewLogger()
	defer logger.Close()

	var colorStart, colorEnd, row string

	if !t.simplified {
		t.stopLoss = t.reached * (1 - t.stopLoss)
	}

	header := fmt.Sprintf(
		"\n| %-10s | %-10s | %-10s | %-10s | %-10s | %-10s | %-10s | %-10s | %-10s |\n",
		"Status", "Symbol", "Current", "Reach", "Reached", "Stop Loss", "Fee", "Target", "Amount",
	)

	if status == BOUGHT {
		colorStart = "\033[1;33m" // Yellow
		colorEnd = "\033[0m"      // Reset
		row = fmt.Sprintf(
			"%s| %-10s | %-10s | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.2f |%s\n",
			colorStart, status, t.symbol, t.current, t.reach, t.reached, t.stopLoss, t.fee, t.target, t.amount, colorEnd,
		)
	} else if status == SOLD {
		colorStart = "\033[1;32m" // Green
		colorEnd = "\033[0m"      // Reset
		row = fmt.Sprintf(
			"%s| %-10s | %-10s | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.2f |%s\n",
			colorStart, status, t.symbol, t.current, t.reach, t.reached, t.stopLoss, t.fee, t.target, t.amount, colorEnd,
		)
	} else {
		row = fmt.Sprintf(
			"| %-10s | %-10s | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.2f |\n",
			status, t.symbol, t.current, t.reach, t.reached, t.stopLoss, t.fee, t.target, t.amount,
		)
	}
	logger.Log(header)
	logger.Log(row)
}

func (t Trade) logError(err error) {
	logger, _ := NewLogger()
	logger.Log("Error: %s", err)
}

func (t *Trade) logStopLoss() {
	logger, _ := NewLogger()
	logger.Log("Stop loss triggered at price %f\n", t.current)
	t.limitLossOrders--
}

func (t *Trade) canBuy() bool {

	if t.simplified {
		return t.simplify(BUY)
	}

	if t.current <= t.reach {
		t.reach = t.current + t.fee + t.target
		t.reached = t.current
		return true
	}

	return false
}

func (t Trade) canSell() bool {

	if t.simplified {
		return t.simplify(SELL) // Simplify trade "STATUS JUST FOR HELPING"
	}

	if t.current <= t.reached*(1-t.stopLoss) {
		t.logStopLoss()
		return true
	}

	return t.current >= (t.reached * (1 + t.target + t.fee))
}

func (t *Trade) simplify(status string) bool {
	buyPrice := getSimplifiedBuyPrice()
	sellPrice := getSimplifiedSellPrice()

	if sellPrice == -1 && buyPrice == -1 {
		t.stopLoss = t.reach * (1 - t.target) * (1 - (t.fee / 2))
	}

	if status == BUY {
		if buyPrice == -1 { // minus one for dynamic price
			t.reached = t.current
			t.reach = t.current + t.fee + t.target
			return true
		}
		if t.current <= buyPrice {
			t.reached = t.current
			return true
		}

	}

	if status == SELL {
		if t.current == t.stopLoss {
			t.logStopLoss()
			return true
		}

		if sellPrice == -1 && t.current == t.reach { // minus one for dynamic price
			t.reached = t.current
			return true
		}
		if sellPrice != -1 && t.current >= sellPrice {
			t.reached = t.current
			return true
		}
	}

	return false
}

func (t *Trade) placeOrder(api *Api, status string) string {

	err := api.PlaceOrderOnApi(status, t.symbol, t.amount, t.current)

	if err != nil {
		t.logError(err)
		if status == BUY {
			return NONE
		} else if status == SELL {
			return BOUGHT
		}
	}
	t.reached = t.current
	if status == BUY {
		status = BOUGHT
	} else if status == SELL {
		status = SOLD
	}

	return status

}

func (t *Trade) refresh(api *Api) {

	media, _ := api.GetMediaFromApi(t.symbol)
	discount := 1 - (t.target + t.fee)

	t.reach = media * discount
	current, _ := api.GetPriceFromApi(t.symbol)
	t.current = current

}
