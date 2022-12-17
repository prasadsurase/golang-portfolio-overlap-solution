package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"geektrust/asset"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
)

/*
type FileData struct {
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

type MutualFund struct {
	Id     int
	Name   string
	Stocks []*Stock
}

// TODO: refactor this to generic function
func getNewStockId(arr *[]*Stock) int {
	max := 0

	for _, obj := range *arr {
		if obj.Id > max {
			max = obj.Id
		}
	}

	return max + 1
}

// TODO: refactor this to generic function
func getNewMutualFundId(arr *[]*MutualFund) int {
	max := 0

	for _, obj := range *arr {
		if obj.Id > max {
			max = obj.Id
		}
	}

	return max + 1
}
*/

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

	var mutualFunds []*asset.MutualFund
	var stocks []*asset.Stock
	var currentPortFolio asset.CurrentPortFolio

	mutualFunds, stocks, err = populateSeedData([]string{"./sample_input/stock_data1.json", "./sample_input/stock_data2.json"})
	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		args := scanner.Text()
		argList := strings.Fields(args)
		fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		switch argList[0] {
		case "CURRENT_PORTFOLIO":
			for i := 1; i <= len(argList)-1; i++ {
				for _, mf := range mutualFunds {
					if argList[i] == mf.Name {
						currentPortFolio.MutualFunds = append(currentPortFolio.MutualFunds, mf)
					}
				}
			}
			fmt.Println(argList)
		case "CALCULATE_OVERLAP":
			fmt.Println(argList)
			var mfs []*asset.MutualFund
			for i := 1; i <= len(argList)-1; i++ {
				for _, mf := range mutualFunds {
					if argList[i] == mf.Name {
						mfs = append(mfs, mf)
					}
				}
			}

			for _, f := range mfs {
				fStocks := mapset.NewSet[asset.Stock]()
				for _, s := range f.Stocks {
					fStocks.Add(*s)
				}
				for _, cmf := range currentPortFolio.MutualFunds {
					// var commonStocks []*Stock
					cmfStocks := mapset.NewSet[asset.Stock]()
					for _, s := range cmf.Stocks {
						cmfStocks.Add(*s)
					}

					commonStocks := fStocks.Intersect(cmfStocks)

					var commonStocksSize float64 = float64(len(commonStocks.ToSlice()))
					var fStockSize float64 = float64(len(f.Stocks))
					var cmfStocksSize float64 = float64(len(cmf.Stocks))
					overlap := ((2 * commonStocksSize) / (fStockSize + cmfStocksSize)) * 100.0
					fmt.Println(f.Name, cmf.Name, overlap)
				}
			}
		case "ADD_STOCK":
			fmt.Println(argList)
			var mf asset.MutualFund
			for _, val := range mutualFunds {
				if val.Name == argList[1] {
					mf = *val
					break
				}
			}
			exists := false
			for _, val := range stocks {
				if val.Name == argList[2] {
					exists = true
				}
			}
			if !exists {
				stk := Stock{Id: getNewStockId(&stocks), Name: argList[2]}
				stocks = append(stocks, &stk)
				mf.Stocks = append(mf.Stocks, &stk)
			}
		default:
		}
	}
}

func populateSeedData(filePaths []string) ([]*asset.MutualFund, []*asset.Stock, error) {
	fundsChan := make(chan asset.Fund, 1000)
	wg := sync.WaitGroup{}
	var mutualFunds []*asset.MutualFund
	var stocks []*asset.Stock
	var err error

	for _, filePath := range filePaths {
		wg.Add(1)

		go func(filePath string) {
			defer wg.Done()

			file, err := os.Open(filePath)
			if err != nil {
				panic(err)
			}

			fmt.Println("Successfully Opened seed file")

			var assetData asset.AssetData

			bytes, _ := ioutil.ReadAll(file)

			err = json.Unmarshal(bytes, &assetData)
			if err != nil {
				panic(err)
			}
			file.Close()

			for _, fund := range assetData.Funds {
				fundsChan <- fund
			}

		}(filePath)
	}

	mutualFunds, stocks, err = asset.Read(fundsChan)

	if err != nil {
		return nil, nil, err
	}
	wg.Wait()

	return mutualFunds, stocks, nil
}
