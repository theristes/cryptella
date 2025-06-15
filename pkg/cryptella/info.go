package cryptella

import (
	"fmt"
	"strconv"
)

func (c *Cryptella) ShowInfo() {

	account := c.api.GetAccountFromApi()
	fmt.Println("Account balances:")
	for _, balance := range account.Balances {
		fr, _ := strconv.ParseFloat(balance.Free, 64)
		if fr == 0 {
			continue
		}
		free, locked := balance.Free, balance.Locked
		if free != "0.00000000" || locked != "0.00000000" {
			fmt.Printf("Asset: %s | Free: %s | Locked: %s\n", balance.Asset, free, locked)
		}
	}

}
