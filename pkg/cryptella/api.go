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

func NewApi() *Api {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")

	return &Api{client: binance.NewClient(apiKey, apiSecret)}
}

func (c Api) GetAccountFromApi() binance.Account {
	account, err := c.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Fatalf("Error fetching account information: %v", err)
	}

	return *account
}

func (c Api) GetPriceFromApi(symbol string) (float64, error) {
	prices, err := c.client.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		log.Printf("Error fetching price data: %s", err)
		return 0, err
	}

	price, err := strconv.ParseFloat(prices[0].Price, 64)
	if err != nil {
		log.Printf("Error parsing price data: %s", err)
		return 0, err
	}

	return price, nil
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

func (c Api) PlaceOrderOnApi(status string, symbol string, amount float64, price float64) error {

	_, err := c.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideType(status)).
		Type(binance.OrderTypeMarket).
		Quantity(fmt.Sprintf("%.8f", amount)).
		Do(context.Background())
	return err
}

func (c Api) GetOrderHistoryFromApi(symbol string) ([]binance.Order, error) {
	orders, err := c.client.NewListOrdersService().Symbol(symbol).Do(context.Background())
	if err != nil {
		log.Printf("Error fetching order history: %s", err)
		return nil, err
	}

	convertedOrders := make([]binance.Order, len(orders))
	for i, order := range orders {

		if order.Price == "0.00000000" {
			// Calculate the average price for market orders
			cummulativeQuoteQty, _ := strconv.ParseFloat(order.CummulativeQuoteQuantity, 64)
			executedQty, _ := strconv.ParseFloat(order.ExecutedQuantity, 64)
			if executedQty > 0 {
				order.Price = fmt.Sprintf("%.8f", cummulativeQuoteQty/executedQty)
			} else {
				order.Price = "0.00000000"
			}
		}
		convertedOrders[i] = *order
	}
	return convertedOrders, nil
}

func (c Api) GetMediaFromApi(symbol string) (float64, error) {
	interval := "1m"
	limit := 60
	values, err := c.GetCandlesFromApi(symbol, interval, limit)
	if err != nil {
		return 0, err
	}

	var sum float64
	for _, value := range values {
		sum += value.Close
	}
	media := sum / float64(len(values))
	return media, nil
}
