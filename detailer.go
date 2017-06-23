package main

import (
	"fmt"
	"os"

	"github.com/MaxSem/detailer/formats"
)

type args struct {
	fileName string
	dryRun   bool
	jsonMode bool
}

func die(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	os.Exit(1)
}

func dieError(err error) {
	die("Error: %s", err.Error())
}

func readArgs() string {
	if len(os.Args) <= 1 {
		die("Usage: detailer [--truncate] <filename>")
	}
	return os.Args[1]
}

func readImage(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		dieError(err)
	}

	stat, err := file.Stat()
	if err != nil {
		dieError(err)
	}

	result, err := formats.ParseJpeg(file)
	if err != nil {
		dieError(err)
	}
	fmt.Printf("JPEG ends at %d bytes out of %d\n", result.DataSize, stat.Size())
}

func main() {
	fileName := readArgs()
	readImage(fileName)
}
