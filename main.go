package main

import (
	"log"
	"os"
	"sync"
	"time"
	"flag"
)

type settings struct {
	OutputFolder    string
	OutputIndexFile string
	MessageLevel    int
	InputSymbolFile string
}

var runtimeSettings settings

func main() {
	log.SetFlags(log.LstdFlags)
	log.SetOutput(os.Stdout)

	flag.StringVar(&runtimeSettings.InputSymbolFile,"symbols","symbols.txt","-symbols=symbols.txt")
	messageLevel := flag.String("level","Warning","-level=[Debug|Info|Warning|Fatal]")
	flag.StringVar(&runtimeSettings.OutputFolder, "outFolder","output", "-outFolder=output")
	flag.StringVar(&runtimeSettings.OutputIndexFile, "outIndex", "_index.idx", "-outindex=index.txt")

	flag.Parse()

	runtimeSettings.MessageLevel = parseLevel(*messageLevel)
	setLevel(runtimeSettings.MessageLevel)

	err := os.MkdirAll(runtimeSettings.OutputFolder, os.ModeDir)
	if err != nil {
		fatal(err, "Could not create output folder", runtimeSettings.OutputFolder)
	}

	symbols := readSymbols(runtimeSettings.InputSymbolFile)

	tickChan := time.Tick(time.Second)

	feederChan := feeder(symbols)

	sinkChan := make(chan stats, 10)

	doneChan := make(chan struct{})

	outputIndexFile := runtimeSettings.OutputFolder + "/" + runtimeSettings.OutputIndexFile

	wg := sync.WaitGroup{}

	wg.Add(1)
	go sink(sinkChan, outputIndexFile, doneChan, &wg)

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
	sinkChan <- stats{"aaa",0,0,0,0}
	close(doneChan)
	sinkChan <- stats{"bbb",0,0,0,0}
	sinkChan <- stats{"ccc",0,0,0,0}

	debug("Symbol channel is now empty. Waiting for pipelines to finish.")
	wg.Wait()
	debug("All done. Exiting...")
}
