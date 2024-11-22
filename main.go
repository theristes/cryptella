package main

import (
	"flag"

	cryptella "cryptella.com/src"
)

var c *cryptella.Cryptella
var info bool

func init() {

	flag.BoolVar(&info, "info", false, "Show information")
	flag.Parse()

	c = cryptella.New()
}

func main() {

	if info {
		c.ShowInfo()
		c.ShowMarketInfo()
		return
	}

	c.Start()
}
