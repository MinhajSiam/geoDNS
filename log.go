package main

import (
	"log"
	"os"
	"time"
)

type LogToFile struct {
	fn      string
	file    *os.File
	closing chan chan error
}

var ltf *LogToFile

func newLogToFile(fn string) *LogToFile {
	return &LogToFile{
		fn:      fn,
		file:    nil,
		closing: make(chan chan error),
	}
}

func logPrintf(format string, a ...interface{}) {
	if *flaglog {
		log.Printf(format, a...)
	}
}

func logPrintln(a ...interface{}) {
	if *flaglog {
		log.Println(a...)
	}
}

func logToFileMonitor() {
	for {
		select {
		case errc := <-ltf.closing:
			if ltf.file != nil {
				log.SetOutput(os.Stderr)
				ltf.file.Close()
				ltf.file = nil
			}
			errc <- nil
			return
		case <-time.After(time.Duration(5 * time.Second)):
			if fi, err := os.Stat(ltf.fn); err != nil || fi.Size() == 0 {
				// it has rotated - first check we can open the new file
				if f, err := os.OpenFile(ltf.fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err == nil {
					log.SetOutput(f)
					log.Printf("Rotating log file")
					ltf.file.Close()
					ltf.file = f
				}
			}
		}
	}
}

func logToFile(fn string) {

	ltf = newLogToFile(fn)

	var err error
	ltf.file, err = os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error writing log file: %v", err)
	}
	// we deliberately do not close logFile here, because we keep it open pretty much for ever

	log.SetOutput(ltf.file)
	log.Printf("Opening log file")

	go logToFileMonitor()
}

func logClose() {
	if ltf != nil {
		log.Printf("Closing log file")
		errc := make(chan error)
		ltf.closing <- errc
		_ = <-errc
		close(ltf.closing)
		ltf = nil
	}
}
