package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

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

	mutualFunds := make([]*MutualFund, 0)
	stocks := make([]*Stock, 0)
	var currentPortFolio CurrentPortFolio

	populateSeedData(&mutualFunds, &stocks)

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
			var mfs []*MutualFund
			for i := 1; i <= len(argList)-1; i++ {
				for _, mf := range mutualFunds {
					if argList[i] == mf.Name {
						mfs = append(mfs, mf)
					}
				}
			}

			for _, f := range mfs {
				fStocks := mapset.NewSet[Stock]()
				for _, s := range f.Stocks {
					fStocks.Add(*s)
				}
				for _, cmf := range currentPortFolio.MutualFunds {
					// var commonStocks []*Stock
					cmfStocks := mapset.NewSet[Stock]()
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
			var mf MutualFund
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

func populateSeedData(mutualFunds *[]*MutualFund, stocks *[]*Stock) {
	seedJSONFile, err := os.Open("./sample_input/stock_data.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened seed file")
	defer seedJSONFile.Close()

	byteValue, _ := ioutil.ReadAll(seedJSONFile)
	var result FileData

	json.Unmarshal([]byte(byteValue), &result)

	for _, v := range result.Funds {
		var mf *MutualFund
		exists := false
		for _, f := range *mutualFunds {
			if f.Name == v.Name {
				exists = true
				mf = f
			}
		}
		if !exists {
			mf = &MutualFund{Id: getNewMutualFundId(mutualFunds), Name: v.Name}
			*mutualFunds = append(*mutualFunds, mf)
		}

		for _, name := range v.Stocks {
			exists := false
			for _, sk := range *stocks {
				if sk.Name == name {
					exists = true
					mf.Stocks = append(mf.Stocks, sk)
				}
			}
			if !exists {
				stk := Stock{Id: getNewStockId(stocks), Name: name}
				*stocks = append(*stocks, &stk)
				mf.Stocks = append(mf.Stocks, &stk)
			}
		}
	}
}
