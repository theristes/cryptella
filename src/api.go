package cryptella

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

func (c *Cryptella) fillApi() {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")

	c.api = Api{client: binance.NewClient(apiKey, apiSecret)}
}

func (c Api) GetAccountFromApi() binance.Account {
	account, err := c.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Fatalf("Error fetching account information: %v", err)
	}

	return *account
}

func (c Api) GetPriceFromApi(symbol string) (float64, float64, float64, float64, error) {

	res, err := c.client.NewKlinesService().Symbol(symbol).Interval("1s").Limit(1).Do(context.Background())
	if err != nil {
		log.Printf("Error fetching price data: %s", err)
		return 0, 0, 0, 0, err
	}

	close, _ := strconv.ParseFloat(res[0].Close, 64)
	open, _ := strconv.ParseFloat(res[0].Open, 64)
	high, _ := strconv.ParseFloat(res[0].High, 64)
	low, _ := strconv.ParseFloat(res[0].Low, 64)

	return close, open, high, low, nil
}

func (c Api) GetCandlesFromApi(symbol string, interval string, limit int) ([]Candle, error) {

	res, err := c.client.NewKlinesService().Symbol(symbol).Interval(interval).Limit(limit).Do(context.Background())
	if err != nil {
		log.Printf("Error fetching candles data: %s", err)
		return nil, err
	}

	candles := make([]Candle, len(res))
	for i, kline := range res {
		open, _ := strconv.ParseFloat(kline.Open, 64)
		close, _ := strconv.ParseFloat(kline.Close, 64)
		high, _ := strconv.ParseFloat(kline.High, 64)
		low, _ := strconv.ParseFloat(kline.Low, 64)
		qav, _ := strconv.ParseFloat(kline.QuoteAssetVolume, 64)
		tbbv, _ := strconv.ParseFloat(kline.TakerBuyBaseAssetVolume, 64)
		tbqav, _ := strconv.ParseFloat(kline.TakerBuyQuoteAssetVolume, 64)

		volume, _ := strconv.ParseFloat(kline.Volume, 64)
		openTime := time.Unix(kline.OpenTime/1000, 0)
		closeTime := time.Unix(kline.CloseTime/1000, 0)

		candles[i] = Candle{
			OpenTime:  openTime,
			Open:      open,
			Close:     close,
			CloseTime: closeTime,
			High:      high,
			Low:       low,
			Volume:    volume,
			QAV:       qav,
			TradeNum:  kline.TradeNum,
			TBBV:      tbbv,
			TBQAV:     tbqav,
		}
	}
	return candles, nil
}

func (c Api) PlaceBuyOrderOnApi(symbol string, amount float64, price float64) error {

	_, err := c.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeMarket).
		Quantity(fmt.Sprintf("%.8f", amount)).
		// Price(fmt.Sprintf("%.8f", price)).
		// Type(binance.OrderTypeLimit).
		// TimeInForce(binance.TimeInForceTypeGTC).
		Do(context.Background())

	if err != nil {
		log.Printf("Error placing buy order: %v", err)
		return err
	}

	return nil
}

func (c Api) PlaceSellOrderOnApi(symbol string, amount float64, price float64) error {

	_, err := c.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeMarket).
		Quantity(fmt.Sprintf("%.8f", amount)).
		// Price(fmt.Sprintf("%.8f", price)).
		// Type(binance.OrderTypeLimit).
		// TimeInForce(binance.TimeInForceTypeGTC).
		Do(context.Background())
	if err != nil {
		log.Printf("Error placing sell order: %v", err)
		return err
	}

	return nil
}

func (c Api) GetBalanceFreeApi(currency string) (float64, error) {

	account, err := c.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return 0, err
	}

	for _, balance := range account.Balances {
		if balance.Asset == currency {
			free, _ := strconv.ParseFloat(balance.Free, 64)
			return free, nil
		}
	}

	return 0, nil
}

func (c Api) GetFiltersFromApi(symbol string) (minQty, maxQty, stepSize, minNotional, minPrice, maxPrice, tickSize, multiplierUp, multiplierDown float64, err error) {
	exchangeInfo, err := c.client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	for _, s := range exchangeInfo.Symbols {
		if s.Symbol == symbol {
			for _, filter := range s.Filters {
				switch filter["filterType"] {
				case "LOT_SIZE":
					minQty, _ = strconv.ParseFloat(filter["minQty"].(string), 64)
					maxQty, _ = strconv.ParseFloat(filter["maxQty"].(string), 64)
					stepSize, _ = strconv.ParseFloat(filter["stepSize"].(string), 64)
				case "MIN_NOTIONAL", "NOTIONAL":
					minNotional, _ = strconv.ParseFloat(filter["minNotional"].(string), 64)
				case "PRICE_FILTER":
					minPrice, _ = strconv.ParseFloat(filter["minPrice"].(string), 64)
					maxPrice, _ = strconv.ParseFloat(filter["maxPrice"].(string), 64)
					tickSize, _ = strconv.ParseFloat(filter["tickSize"].(string), 64)
				case "PERCENT_PRICE":
					multiplierUp, _ = strconv.ParseFloat(filter["multiplierUp"].(string), 64)
					multiplierDown, _ = strconv.ParseFloat(filter["multiplierDown"].(string), 64)
				}
			}

			return minQty, maxQty, stepSize, minNotional, minPrice, maxPrice, tickSize, multiplierUp, multiplierDown, nil
		}
	}
	return 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("filters not found for symbol %s", symbol)
}
