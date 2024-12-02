# Cryptella

Cryptella is a cryptocurrency trading bot that interacts with the Binance API to automate trading strategies. It fetches real-time price data, places buy and sell orders, and logs trading activities.

## Features

- Fetches real-time price data from Binance
- Places buy and sell orders based on predefined strategies
- Logs trading activities in a structured format
- Supports simplified trading mode
- Configurable trading parameters via environment variables

## Installation

1. Clone the repository:

```sh
git clone https://github.com/theristes/cryptella.git
cd cryptella
```

2. Install dependencies:

```sh
go mod tidy
```

3. Set up your environment variables in a .env file:

```sh
# Binance API Keys
BINANCE_API_KEY=your_api_key
BINANCE_API_SECRET=your_api_secret

# Trading parameters
SYMBOL=XLMUSDT                 # The cryptocurrency pair you're trading
AMOUNT=15                      # The amount of the asset to use for trading
FEE=0.002                      # The fee
TARGET=0.002                   # The target profit
LIMIT_LOSS_ORDERS=5            # The number of loss orders to place
SIMPLIFIED=true                # Simplified mode
BUY_PRICE=-1
SELL_PRICE=-1
STOP_LOSS=0.01                 # The stop loss
```

# Usage

Run the trading bot:

```sh
go run main.go
```

To show information:

```sh
go run main.go --info
```

To show order history:

```sh
go run main.go --history
```

# Configuration

You can configure the trading bot using environment variables. Here are the available options:

```properties
BINANCE_API_KEY: Your Binance API key.
BINANCE_API_SECRET: Your Binance API secret
SYMBOL: The cryptocurrency pair you are trading (e.g., XLMUSDT)
AMOUNT: The amount of the asset to use for trading
FEE: The trading fee (e.g., 0.002 for 0.2%)
TARGET: The target profit (e.g., 0.002 for 0.2%)
LIMIT_LOSS_ORDERS: The number of loss orders to place
SIMPLIFIED: Enable simplified trading mode (true or false)
BUY_PRICE: The buy price for simplified mode (-1 for dynamic price)
SELL_PRICE: The sell price for simplified mode (-1 for dynamic price)
STOP_LOSS: The stop loss percentage (e.g., 0.01 for 1%)
```
