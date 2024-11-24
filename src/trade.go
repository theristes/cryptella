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

	discount := 1 - (target + fee)
	media, _ := api.GetMediaFromApi(symbol)
	reach := media * discount
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
	}
}

func (t Trade) canTrade() bool {
	return t.limitLossOrders > 0
}

func (t Trade) logInfo(status string) {
	logger, _ := NewLogger()
	defer logger.Close()

	header := fmt.Sprintf(
		"\n| %-10s | %-10s | %-10s | %-10s | %-10s | %-10s | %-10s | %-10s | %-10s |\n",
		"Status", "Symbol", "Current", "Reach", "Reached", "Stop Loss", "Fee", "Target", "Amount",
	)

	row := fmt.Sprintf(
		"| %-10s | %-10s | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.4f | %-10.2f |\n",
		status, t.symbol, t.current, t.reach, t.reached, t.stopLoss, t.fee, t.target, t.amount,
	)
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

func (t Trade) canBuy() bool {
	t.reach = t.reach + 0.001
	return t.current <= t.reach
}

func (t Trade) canSell() bool {

	if t.current >= t.reached*(1-t.stopLoss) {
		t.logStopLoss()
		return true

	}
	return t.current >= t.getSellPrice(t.reached)
}

func (t Trade) getSellPrice(price float64) float64 {
	return price * (1 + t.target + t.fee)
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
