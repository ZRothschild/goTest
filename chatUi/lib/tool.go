package lib

import (
	"log"
	"runtime"
)

func Log(err error, msg string) {
	_, file, line, ok := runtime.Caller(1)
	if err != nil && ok {
		log.Fatalf("%s: %s  file name %s line %d", msg, err, file, line)
	}
}
