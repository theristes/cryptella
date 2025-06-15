package cryptella

import (
	"fmt"
	"log"
)

func (c *Cryptella) ShowOrderHistory() {
	orders, err := c.api.GetOrderHistoryFromApi(c.trade.symbol)
	if err != nil {
		log.Fatalf("Error fetching order history: %v", err)
	}

	for _, order := range orders {
		fmt.Printf("Order ID: %d, Symbol: %s, Price: %s, Quantity: %s, Status: %s\n",
			order.OrderID, order.Symbol, order.Price, order.OrigQuantity, order.Status)
	}
}
