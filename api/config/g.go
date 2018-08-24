package config

import (
	"log"
)

// change log:
const (
	VERSION = "0.0.1"
)

func init() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
