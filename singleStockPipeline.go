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
	debug("Downloading data for", symbol)
	defer wg.Done()

	content, err := downloadData(symbol)
	if err != nil || content == "" {
		// skip the rest of the pipeline
		return
	}

	debug("Successfully got data for", symbol, ":", content)

	stats, err := parseContent(content)
	if err != nil {
		// skip the rest of pipeline
		return
	}

	stats.symbol = symbol

	sinkChan <- stats
}

var regex = regexp.MustCompile("^[a0-9]+,([\\dd.]+),([\\dd.]+)$")

func parseContent(content string) (stats, error) {

	lines := strings.Split(content, "\n")
	min := math.MaxFloat64
	max := float64(0.0)
	price := float64(0.0)

	for _, line := range lines {
		debug("parsing line", line)
		match := regex.FindAllStringSubmatch(line, -1)
		if match == nil || len(match[0]) != 3 {
			continue
		}

		debug("parsed line", match)

		daymin, err := strconv.ParseFloat(match[0][2], 64)
		if err != nil {
			warning(err, "Error when parsing float", match[0])
		} else {
			if daymin < 0.001 {
				warning("encountered an abnormal daily min value", daymin)
			} else {
				min = math.Min(min, daymin)
			}
		}

		daymax, err := strconv.ParseFloat(match[0][1], 64)
		if err != nil {
			warning(err, "Error when parsing float", match[1])
		} else {
			max = math.Max(max, daymax)

			// treat the day's max price as the current price
			price = daymax
		}

		debug("parsed min:", daymin, "max:", daymax)
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
