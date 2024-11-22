package cryptella

import (
	"log"
	"os"
	"strconv"
	"time"

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
	tradeMode := os.Getenv("TRADE_MODE")

	interval := os.Getenv("INTERVAL")
	maxTrades, _ := strconv.Atoi(os.Getenv("MAX_TRADES"))
	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))
	buyThresehold, _ := strconv.ParseFloat(os.Getenv("BUY_THRESEHOLD"), 64)
	balanceAsset, _ := strconv.ParseFloat(os.Getenv("BALANCE_ASSET"), 64)

	cryptella := &Cryptella{
		symbol:        symbol,
		fee:           fee,
		target:        target,
		maxTrades:     maxTrades,
		interval:      interval,
		limit:         limit,
		buyThresehold: buyThresehold,
		balanceAsset:  balanceAsset,
		tradeMode:     tradeMode,
		status:        NONE,
	}

	cryptella.fillApi()
	cryptella.fillAmount()
	cryptella.fillLowPrice()

	return cryptella
}

func (c *Cryptella) Start() {

	i := 0

	for {

		if c.maxTrades == 0 {
			log.Println("Max trades reached")
			break
		}

		c.analyze()

		if c.status == BUY {
			i = 0
			c.buy()
		}

		if c.status == SELL {
			i = 0
			c.sell()
		}

		i++

		if i == 100 {
			c.fillLowPrice()
			i = 0
		}

		time.Sleep(500 * time.Millisecond)

	}

}
