package main

import (
	"log"
)

const (
	Debug = iota
	Info
	Warning
	Error
	Fatal
)

var currentLevel = Warning

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
