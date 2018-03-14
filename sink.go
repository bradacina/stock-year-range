package main

import (
	"os"
	"fmt"
	"sync"
)

func writeToIndex(file *os.File, thestats stats) {
		debug("going to write to file", fmt.Sprintf("%v", thestats))
		_, err := file.WriteString(
			fmt.Sprintf("%v,%v,%v,%v,%v\n", 
			thestats.symbol,
			thestats.price,
			thestats.min,
			thestats.max,
			thestats.percentage))
		if err != nil {
			fatal(err, "Cannot write to output index file")
		}
	}

func sink(
	sinkChan <-chan stats,
	outputFile string, 
	done <-chan struct{},
	wg *sync.WaitGroup) {
		
		file, err := os.Create(outputFile)
		if err != nil {
			fatal(err, "Cannot create output file", outputFile)
		}
		defer file.Close()

	isDone := false
	for {
		select{
		case stats := <- sinkChan:
			writeToIndex(file,stats)
			break
		case <- done:
			isDone = true
			done = nil
			break
		default:
			if isDone {
				wg.Done()
				return
			}
		}
	}
}