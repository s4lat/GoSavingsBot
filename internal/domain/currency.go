package domain

type Currency struct {
	Name   string
	Ticker Ticker
	Symbol string
}

type Ticker = string

var (
	currencyMap = map[Ticker]Currency{
		"USD": {
			Name:   "USD, United States Dollar",
			Ticker: "USD",
			Symbol: "$",
		},
		"RUB": {
			Name:   "RUB, Russian Ruble",
			Ticker: "RUB",
			Symbol: "₽",
		},
	}
)

func GetCurrency(ticker Ticker) (Currency, bool) {
	curr, ok := currencyMap[ticker]
	return curr, ok
}
