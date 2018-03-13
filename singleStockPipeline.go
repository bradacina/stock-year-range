package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"regexp"
	"strings"
)

type stats struct {
	symbol string
	price float32
	min float32
	max float32
	percentage float32
}

func singleStockPipeline(
	symbol string, 
	wg *sync.WaitGroup,
	sinkChan chan<- stats) {
	debug("Downloading data for", symbol)
	defer wg.Done()

	content, err := downloadData(symbol)
	if err != nil || content == ""{
		// skip the rest of the pipeline
		return
	}

	debug("Successfully got data for", symbol, ":", content)

	stats, err := parseContent(content)
	if err != nil {
		// skip the rest of pipeline
		return
	}

	sinkChan <- stats
}

var regex = regexp.MustCompile("^[a0-9]+,([\\dd.]+),([\\dd.]+)$")
func parseContent(content string) (stats, error) {
	
	lines := strings.Split(content, "\n")
	//list := make([]stats, len(lines))

	for _, line := range lines {
		debug("parsing line", line)
		match := regex.FindAllStringSubmatch(line, -1)
		if match == nil {
			continue
		}

		debug("parsed line", match)
	}
	
	return stats{
		symbol:"haha",
		price:123,
		min:100,
		max:200,
		percentage : .2}, nil
}

func downloadData(symbol string) (string, error) {

	outputFile := runtimeSettings.OutputFolder + "/" + symbol + ".txt"
	stat, _ := os.Stat(outputFile)
	if stat != nil {
		info("Found", outputFile, ". Skipping http download.")

		fileBytes, err := ioutil.ReadFile(outputFile)
		if err != nil {
			warning("Could not read contents of", outputFile, ". Skipping.")
			return "", err
		}

		return string(fileBytes), nil
	}
	url := fmt.Sprintf("https://www.google.com/finance/getprices?q=%s&p=1Y&f=d,h,l&i=86401", symbol)
	debug("Downloading data from", url)

	resp, err := http.Get(url)
	if err != nil {
		warning(err, "Error when retrieving data over http. Skipping.")
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		warning("Server returned", resp.StatusCode, "status code. Skipping...")
		return "", http.ErrAbortHandler
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		warning(err, "Could not read response body. Skipping...")
		return "", err
	}

	debug("Saving data to", outputFile)

	err = ioutil.WriteFile(outputFile, bodyBytes, 0)
	if err != nil {
		warning(err, "Error writing output file", outputFile)
	}

	return string(bodyBytes), nil
}
