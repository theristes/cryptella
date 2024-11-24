package cryptella

import (
	"log"
	"time"

	"github.com/joho/godotenv"
)

func NewCryptella() *Cryptella {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	api := NewApi()
	trade := NewTrade(api)

	cryptella := &Cryptella{

		api:    api,
		trade:  trade,
		status: NONE,
		stopCh: make(chan struct{}),
	}

	return cryptella
}

func (c Cryptella) listenPrice() {
	ticker := time.NewTicker(700 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			price, err := c.api.GetPriceFromApi(c.trade.symbol)
			if err != nil {
				log.Printf("Error fetching price: %v", err)
				continue
			}
			
			if price != c.trade.current {
				c.analyze(price)
			}
		case <-c.stopCh:
			return
		}
	}
}

func (c Cryptella) Start() {
	c.listenPrice()
}

func (c Cryptella) Stop() {
	c.stopCh <- struct{}{}
}
