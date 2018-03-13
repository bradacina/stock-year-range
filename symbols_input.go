package main

import (
	"io/ioutil"
	"strings"
)

func readSymbols(filename string) []string {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fatal(err, "Could not read", filename)
	}

	lines := strings.Split(string(fileBytes), "\r\n")

	debug("Parsed the following input symbols", lines)

	return lines
}
