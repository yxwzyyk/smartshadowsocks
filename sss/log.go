package sss

import (
	"log"
	"os"
)

type SocksLog struct {
	Debug bool
}

var SLog SocksLog

var info = log.New(os.Stdout, "[INFO] ", log.Ltime)
var debug = log.New(os.Stdout, "[DEBUG] ", log.Ltime)

func (d *SocksLog) Printf(format string, args ...interface{}) {
	info.Printf(format, args...)
}

func (d *SocksLog) Println(args ...interface{}) {
	info.Println(args...)
}

func (d *SocksLog) DebugPrintf(format string, args ...interface{}) {
	if d.Debug {
		debug.Printf(format, args...)
	}
}

func (d *SocksLog) DebugPrintln(args ...interface{}) {
	if d.Debug {
		debug.Println(args...)
	}
}
