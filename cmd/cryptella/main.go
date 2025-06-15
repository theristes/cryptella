package main

import (
	"flag"

	"cryptella.com/pkg/cryptella"
)

var c *cryptella.Cryptella
var info bool
var history bool

func init() {

	flag.BoolVar(&info, "info", false, "Show information")
	flag.BoolVar(&history, "history", false, "Show order history")

	flag.Parse()

	c = cryptella.NewCryptella()
}

func main() {

	if info {
		c.ShowInfo()
		c.ShowMarketInfo()
		return
	}

	if history {
		c.ShowOrderHistory()
		return
	}

	c.Start()
}
