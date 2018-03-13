package main

import (
	"os"
	"sync"
	"time"
)

type settings struct {
	OutputFolder string
}

var runtimeSettings settings

func main() {
	setLevel(Debug)

	runtimeSettings.OutputFolder = "output"

	os.MkdirAll(runtimeSettings.OutputFolder, os.ModeDir)

	tickChan := time.Tick(time.Second)

	feederChan := feeder([]string{"AAPL", "vvv", "ccc"})

	wg := sync.WaitGroup{}

tickerLoop:
	for {
		select {
		case <-tickChan:
			symbol, ok := <-feederChan
			if !ok {
				// we don't have any more values in the feeder
				break tickerLoop
			}

			wg.Add(1)
			go singleStockPipeline(symbol, &wg)
			break
		}
	}

	debug("Symbol channel is now empty. Waiting for pipelines to finish.")
	wg.Wait()
	debug("All done. Exiting...")
}
