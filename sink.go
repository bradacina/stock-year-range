package main

import (
	"os"
	"fmt"
)

type sinkDelegate (func(*os.File, stats))

func fileSinkDelegate(file *os.File, thestats stats) {
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
	f sinkDelegate, 
	done <-chan struct{}) {
		
		file, err := os.Create(outputFile)
		if err != nil {
			fatal(err, "Cannot create output file", outputFile)
		}
		defer file.Close()

	for {
		select{
		case stats := <- sinkChan:
			f(file,stats)
			break
		case <- done:
			return
		}
	}
}