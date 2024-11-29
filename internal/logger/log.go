package logger

import (
	"io"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func DisableLog() {
	log.SetOutput(io.Discard)
}

func Error(err error) {
	log.Println("ERROR", err.Error())
}

func Info(message string) {
	log.Println("INFO", message)
}
