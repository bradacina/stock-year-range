package main

import (
	"log"
	"strings"
)

const (
	Debug = iota
	Info
	Warning
	Error
	Fatal
)

var currentLevel = Warning

func parseLevel(lvl string) int {
	lvl = strings.ToLower(lvl)

	switch(lvl){
	case "debug": return Debug
	case "info": return Info
	case "warning": return Warning
	case "fatal": return Fatal
	default: return Warning
	}
}

func setLevel(level int) {
	currentLevel = level
}

func debug(v ...interface{}) {
	if currentLevel <= Debug {
		log.Println("[DEBUG]", v)
	}
}

func info(v ...interface{}) {
	if currentLevel <= Info {
		log.Println("[INFO]", v)
	}
}

func warning(v ...interface{}) {
	if currentLevel <= Warning {
		log.Println("[WARNING]", v)
	}
}

func fatal(v ...interface{}) {
	if currentLevel <= Fatal {
		log.Fatalln("[FATAL]", v)
	}
}
