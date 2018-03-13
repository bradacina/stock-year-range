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

	err := os.MkdirAll(runtimeSettings.OutputFolder, os.ModeDir)
	if err != nil {
		fatal(err, "Could not create output folder", runtimeSettings.OutputFolder)
	}

	tickChan := time.Tick(time.Second)

	feederChan := feeder([]string{"NUH.AX"})

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
