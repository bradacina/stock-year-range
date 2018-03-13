package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type stats struct {
	symbol     string
	price      float64
	min        float64
	max        float64
	percentage float64
}

func singleStockPipeline(
	symbol string,
	wg *sync.WaitGroup,
	sinkChan chan<- stats) {
	debug(symbol, ":Downloading data")
	defer wg.Done()

	content, err := downloadData(symbol)
	if err != nil || content == "" {
		// skip the rest of the pipeline
		return
	}

	debug(symbol, ":Successfully got data:", content)

	stats, err := parseContent(symbol, content)
	if err != nil {
		// skip the rest of pipeline
		return
	}

	stats.symbol = symbol

	sinkChan <- stats
}

var regex = regexp.MustCompile("^[a0-9]+,([\\dd.]+),([\\dd.]+)$")

func parseContent(symbol, content string) (stats, error) {

	lines := strings.Split(content, "\n")
	min := math.MaxFloat64
	max := float64(0.0)
	price := float64(0.0)

	for _, line := range lines {
		debug(symbol, ":parsing line in ", line)
		match := regex.FindAllStringSubmatch(line, -1)
		if match == nil || len(match[0]) != 3 {
			continue
		}

		debug(symbol, ":parsed line", match)

		daymin, err := strconv.ParseFloat(match[0][2], 64)
		if err != nil {
			warning(err, symbol, ":Error when parsing float", match[0])
		} else {
			if daymin < 0.001 {
				warning(symbol, ":encountered an abnormal daily min value", daymin)
			} else {
				min = math.Min(min, daymin)
			}
		}

		daymax, err := strconv.ParseFloat(match[0][1], 64)
		if err != nil {
			warning(err, symbol, ":Error when parsing float", match[1])
		} else {
			if daymax < 0.001 {
				warning(symbol, ":encountered an abnormal daily max value", daymax)
			} else {
				max = math.Max(max, daymax)

				// treat the day's max price as the current price
				price = daymax
			}
		}

		debug(symbol, ":parsed min:", daymin, "max:", daymax)
	}

	percentage := (price - min) / (max - min)

	return stats{
		price:      price,
		min:        min,
		max:        max,
		percentage: percentage}, nil
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
	debug(symbol, ":Downloading data from", url)

	resp, err := http.Get(url)
	if err != nil {
		warning(err, symbol, ":Error when retrieving data over http. Skipping.")
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		warning(symbol, ":Server returned", resp.StatusCode, "status code. Skipping...")
		return "", http.ErrAbortHandler
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		warning(err, symbol, ":Could not read response body. Skipping...")
		return "", err
	}

	debug(symbol, ":Saving data to", outputFile)

	err = ioutil.WriteFile(outputFile, bodyBytes, 0)
	if err != nil {
		warning(err, symbol, ":Error writing output file", outputFile)
	}

	return string(bodyBytes), nil
}
