package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MaxSem/detailer/formats"
)

var jsonMode = false

type args struct {
	fileName string
	truncate bool
}

func die(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func dieError(err error) {
	die("Error: %s", err.Error())
}

func readArgs() *args {
	var result args
	flag.BoolVar(&result.truncate, "truncate", false, "Truncate the file")
	flag.BoolVar(&jsonMode, "json", false, "Output results as machine-readable JSON")
	flag.Parse()
	result.fileName = flag.Arg(0)
	if result.fileName == "" {
		fmt.Fprintf(os.Stderr, "Detailer removes junk past the end of media files.\n"+
			"Usage: detailer [-truncate] [-json] <filename>")
		flag.PrintDefaults()
		os.Exit(1)
	}
	return &result
}

func processImage(fileName string, truncate bool) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0777)
	if err != nil {
		dieError(err)
	}

	parser, err := formats.GetParser(fileName)
	if err != nil {
		dieError(err)
	}

	stat, err := file.Stat()
	if err != nil {
		dieError(err)
	}
	origSize := stat.Size()

	result, err := parser.Parse(file)
	if err != nil {
		dieError(err)
	}

	truncated := false
	if truncate && origSize > result.DataSize {
		err = file.Truncate(result.DataSize)
		if err != nil {
			dieError(err)
		}
		truncated = true
	}

	if jsonMode {
		fmt.Printf("{\n\t\"format\": \"%s\",\n\t\"size\": %d,\n\t\"data_size\": %d,\n\t\"truncated\": %t\n}", parser.Format(), origSize, result.DataSize, truncated)
	}
}

func main() {
	params := readArgs()
	processImage(params.fileName, params.truncate)
}
