package cryptella

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func New() *Cryptella {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	symbol := os.Getenv("SYMBOL")
	fee, _ := strconv.ParseFloat(os.Getenv("FEE"), 64)
	target, _ := strconv.ParseFloat(os.Getenv("TARGET"), 64)
	stopLoss, _ := strconv.ParseFloat(os.Getenv("STOP_LOSS"), 64)

	interval := os.Getenv("INTERVAL")
	maxTrades, _ := strconv.Atoi(os.Getenv("MAX_TRADES"))
	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))
	balanceAsset, _ := strconv.ParseFloat(os.Getenv("BALANCE_ASSET"), 64)

	cryptella := &Cryptella{
		symbol:       symbol,
		fee:          fee,
		stopLoss:     stopLoss,
		target:       target,
		maxTrades:    maxTrades,
		interval:     interval,
		limit:        limit,
		balanceAsset: balanceAsset,
		status:       NONE,
	}

	cryptella.fillApi()
	cryptella.fillAmount()
	cryptella.fillBuyThreshold()
	cryptella.fillLowPrice()

	return cryptella
}
func (c *Cryptella) fillBuyThreshold() {
	buyThresehold, _ := strconv.ParseFloat(os.Getenv("BUY_THRESEHOLD"), 64)
	c.buyThresehold = buyThresehold
}

func (c *Cryptella) Start() {

	for {
		if c.maxTrades == 0 {
			fmt.Println("Max trades reached")
			break
		}
		c.analyze()
	}

}
