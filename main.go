package main

import (
	"os"
	"sync"
	"time"
)

type settings struct {
	OutputFolder string
	OutputIndex  string
	MessageLevel int
}

var runtimeSettings settings

func main() {
	runtimeSettings.MessageLevel = Warning
	runtimeSettings.OutputFolder = "output"
	runtimeSettings.OutputIndex = "index.idx"

	setLevel(runtimeSettings.MessageLevel)

	err := os.MkdirAll(runtimeSettings.OutputFolder, os.ModeDir)
	if err != nil {
		fatal(err, "Could not create output folder", runtimeSettings.OutputFolder)
	}

	tickChan := time.Tick(time.Second)

	feederChan := feeder([]string{"NUH.AX", "Z1P"})

	sinkChan := make(chan stats, 10)

	doneChan := make(chan struct{})

	outputIndexFile := runtimeSettings.OutputFolder + "/" + runtimeSettings.OutputIndex
	go sink(sinkChan, outputIndexFile, fileSinkDelegate, doneChan)

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
			go singleStockPipeline(symbol, &wg, sinkChan)
			break
		}
	}

	debug("Symbol channel is now empty. Waiting for pipelines to finish.")
	wg.Wait()
	debug("All done. Exiting...")
}
