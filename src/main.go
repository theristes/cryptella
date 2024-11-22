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

	interval := os.Getenv("INTERVAL")
	maxTrades, _ := strconv.Atoi(os.Getenv("MAX_TRADES"))
	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))
	buyThresehold, _ := strconv.ParseFloat(os.Getenv("BUY_THRESEHOLD"), 64)

	cryptella := &Cryptella{
		symbol:        symbol,
		fee:           fee,
		target:        target,
		maxTrades:     maxTrades,
		interval:      interval,
		buyThresehold: buyThresehold,
		limit:         limit,
		status:        NONE,
	}

	cryptella.fillApi()
	cryptella.fillAmount()
	cryptella.fillLowPrice()

	return cryptella
}

func (c *Cryptella) Start() {

	for {

		if c.maxTrades == 0 {
			log.Println("Max trades reached")
			break
		}

		c.analyze()

		if c.status == BUY {
			c.buy()
		}

		if c.status == SELL {
			c.sell()
		}

		time.Sleep(500 * time.Millisecond)

	}

}
