package main

import (
	"os"
	"bufio"
	"fmt"
)

type sinkDelegate (func(stats))

func newFileSink(outputFile string) sinkDelegate {
	f, err := os.Create(outputFile)
	if err != nil {
		fatal(err, "Cannot create output file", outputFile)
	}

	bufWriter := bufio.NewWriter(f)

	return func(thestats stats) {
		debug("going to write to file %v", thestats)
		_, err := bufWriter.WriteString(fmt.Sprintf("%v", thestats))
		if err != nil {
			fatal(err, "Cannot write to output index file")
		}
	}
}

func sink(sinkChan <-chan stats, f sinkDelegate, done <-chan struct{}) {
	for {
		select{
		case stats := <- sinkChan:
			f(stats)
			break
		case <- done:
			return
		}
	}
}