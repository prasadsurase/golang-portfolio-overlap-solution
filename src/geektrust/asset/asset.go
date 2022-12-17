package asset

type AssetData struct {
	Funds []Fund `json:"funds"`
}

type Fund struct {
	Name   string   `json:"name"`
	Stocks []string `json:"stocks"`
}

type CurrentPortFolio struct {
	MutualFunds []*MutualFund
}

type Stock struct {
	Id   int
	Name string
}

func NewStockID(stocks []*Stock) int {
	max := 0

	for _, stock := range stocks {
		if stock.Id > max {
			max = stock.Id
		}
	}

	return max + 1

}

type MutualFund struct {
	Id     int
	Name   string
	Stocks []*Stock
}

func NewMutualFundID(mutualFunds []*MutualFund) int {
	max := 0

	for _, mutualFund := range mutualFunds {
		if mutualFund.Id > max {
			max = mutualFund.Id
		}
	}

	return max + 1
}

func Read(funds chan Fund) ([]*MutualFund, []*Stock, error) {
	var mutualFunds []*MutualFund
	var stocks []*Stock

	mutualFundsMap := make(map[string]*MutualFund)
	stocksMap := make(map[string]*Stock)

	for fund := range funds {
		if _, ok := stocksMap[fund.Name]; !ok {
			for _, stock := range fund.Stocks {
				stocksMap[stock] = &Stock{Id: 1, Name: stock}
			}
		}

		if _, ok := mutualFundsMap[fund.Name]; !ok {
			mutualFund := &MutualFund{Id: 1, Name: fund.Name}
			for _, stock := range fund.Stocks {
				mutualFund.Stocks = append(mutualFund.Stocks, stocksMap[stock])
			}
			mutualFundsMap[fund.Name] = mutualFund
		}
	}

	for _, mutualFund := range mutualFundsMap {
		mutualFunds = append(mutualFunds, mutualFund)
	}

	for _, stock := range stocksMap {
		stocks = append(stocks, stock)
	}

	return mutualFunds, stocks, nil
}
