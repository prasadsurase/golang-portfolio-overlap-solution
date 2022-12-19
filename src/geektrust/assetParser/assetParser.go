package assetParser

import (
	"encoding/json"
	"fmt"
	"geektrust/asset"
	"io/ioutil"
	"os"
	"sync"
)

func ParseFile(filePath string, fundsChan chan asset.Fund, wg *sync.WaitGroup) error {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully opened seed file")
	var parsedData asset.ParsedData
	bytes, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(bytes, &parsedData)
	if err != nil {
		panic(err)
	}

	defer wg.Done()

	for _, fund := range parsedData.Funds {
		fundsChan <- *fund
	}

	return nil
}
