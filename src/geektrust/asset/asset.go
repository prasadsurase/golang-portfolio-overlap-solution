package asset

import (
	"fmt"
)

type Fund struct {
	Name   string   `json:"name"`
	Stocks []string `json:"stocks"`
}

type ParsedData struct {
	Funds []*Fund
}

type Stock struct {
	Id   int
	Name string
}

type MutualFund struct {
	Id     int
	Name   string
	Stocks []*Stock
}

func ParseData(fund Fund, mutualFundsMap map[string]*MutualFund, stocksMap map[string]*Stock) error {
	fmt.Println("Received Fund:", fund.Name, "Stocks count", len(fund.Stocks))

	if _, ok := stocksMap[fund.Name]; !ok {
		for _, stock := range fund.Stocks {
			stocksMap[stock] = &Stock{Id: 1, Name: stock}
		}
	}

	if _, ok := mutualFundsMap[fund.Name]; !ok {
		mutualFund := &MutualFund{Id: 1, Name: fund.Name}
		for _, stockName := range fund.Stocks {
			mutualFund.Stocks = append(mutualFund.Stocks, stocksMap[stockName])
		}
		mutualFundsMap[fund.Name] = mutualFund
	}
	return nil
}
