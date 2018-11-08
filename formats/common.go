package formats

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type DecodingResult struct {
	DataSize int64
}

type Parser interface {
	Format() string
	Parse (*os.File) (*DecodingResult, error)
	Validate (*os.File) (bool, error)
}

var extensionList = map[string]string{
	"jpeg": "JPEG",
	"jpg":  "JPEG",
	"jpe":  "JPEG",
	"png":  "PNG",
}

var parserList = map[string]Parser{
	"JPEG": formatJpeg{},
	"PNG":  formatPng{},
}


// GetParser returns the function to parse the given file
// The parser is determined from file extension only and
// parsing will fail if file contents mismatch it.
func GetParser(fileName string) (Parser, error) {
	ext := strings.ToLower(path.Ext(fileName))
	if ext != "" && ext[0] == '.' {
		ext = ext[1:]
	}
	format, ok := extensionList[ext]
	if !ok {
		return nil, fmt.Errorf("unrecognized extension '%s'", ext)
	}
	parser, _ := parserList[format]
	return parser, nil
}
