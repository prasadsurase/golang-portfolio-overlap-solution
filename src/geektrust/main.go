package main

import (
	"bufio"
	"fmt"
	"geektrust/asset"
	"geektrust/assetParser"
	"os"
	"strings"
	"sync"
)

type CurrentPortfolio struct {
	MutualFunds []*asset.MutualFund
}

func (cp CurrentPortfolio) CalculateOverlap(mf *asset.MutualFund) {
	for _, cpfund := range cp.MutualFunds {
		var commonStocksSize float64

		for _, cs := range cpfund.Stocks {
			for _, fs := range mf.Stocks {
				if cs.Name == fs.Name {
					commonStocksSize += 1.0
				}
			}
		}

		var fStockSize float64 = float64(len(cpfund.Stocks))
		var cmfStocksSize float64 = float64(len(mf.Stocks))
		overlap := ((2 * commonStocksSize) / (fStockSize + cmfStocksSize)) * 100.0
		fmt.Println(mf.Name, cpfund.Name, overlap)
	}
}

func main() {
	cliArgs := os.Args[1:]

	if len(cliArgs) == 0 {
		fmt.Println("Please provide the input file path")

		return
	}

	filePath := cliArgs[0]
	file, err := os.Open(filePath)

	if err != nil {
		fmt.Println("Error opening the input file")

		return
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	// var mutualFunds []*asset.MutualFund
	// var stocks []*asset.Stock
	var currentPortFolio CurrentPortfolio

	mutualFundsMap := make(map[string]*asset.MutualFund)
	stocksMap := make(map[string]*asset.Stock)
	fundsChan := make(chan asset.Fund, 1000)
	wg1 := sync.WaitGroup{}

	seedFiles := []string{
		"./sample_input/stock_data1.json",
		"./sample_input/stock_data2.json",
		"./sample_input/stock_data3.json",
		"./sample_input/stock_data4.json",
	}

	for _, filePath := range seedFiles {
		wg1.Add(1)
		go assetParser.ParseFile(filePath, fundsChan, &wg1)
	}
	wg1.Wait()

	for len(fundsChan) > 0 {
		asset.ParseData(<-fundsChan, mutualFundsMap, stocksMap)
	}
	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		args := scanner.Text()
		argList := strings.Fields(args)
		fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		fmt.Println(argList)
		switch argList[0] {
		case "CURRENT_PORTFOLIO":
			for i := 1; i <= len(argList)-1; i++ {
				fmt.Println("MF: ", argList[i])
				if _, ok := mutualFundsMap[argList[i]]; ok {
					currentPortFolio.MutualFunds = append(currentPortFolio.MutualFunds, mutualFundsMap[argList[i]])
				}
			}
		case "CALCULATE_OVERLAP":
			for i := 1; i <= len(argList)-1; i++ {
				if _, ok := mutualFundsMap[argList[i]]; ok {
					currentPortFolio.CalculateOverlap(mutualFundsMap[argList[i]])
				} else {
					fmt.Println("FUND NOT FOUND", argList[i])
				}
			}
		case "ADD_STOCK":
			var stk asset.Stock

			if _, ok := stocksMap[argList[2]]; ok {
				stk = *stocksMap[argList[2]]
			} else {
				stk = asset.Stock{Name: argList[2]}
				stocksMap[argList[2]] = &stk
			}
			if _, ok := mutualFundsMap[argList[1]]; ok {
				mf := mutualFundsMap[argList[1]]
				mf.Stocks = append(mf.Stocks, &stk)
				mutualFundsMap[mf.Name] = mf
			} else {
				fmt.Println("FUND NOT FOUND", argList[1])
			}
		}
	}
}
